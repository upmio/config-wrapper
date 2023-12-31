package template

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/kelseyhightower/memkv"

	"github.com/upmio/config-wrapper/app/config/confd/backends"
	"github.com/upmio/config-wrapper/pkg/util"
)

type Config struct {
	StoreClient  backends.StoreClient
	TemplateFile string
	DestFile     string
}

const (
	uid      = 1001
	gid      = 1001
	fileMode = 0644
	dirMode  = 0755
)

// TemplateResource is the representation of a parsed template resource.
type TemplateResource struct {
	Dest        string
	FileMode    os.FileMode
	Gid         int
	Mode        string
	Src         string
	StageFile   *os.File
	Uid         int
	funcMap     map[string]interface{}
	store       memkv.Store
	storeClient backends.StoreClient
}

// NewTemplateResource creates a TemplateResource.
func NewTemplateResource(config Config) (*TemplateResource, error) {
	if config.StoreClient == nil {
		return nil, errors.New("A valid StoreClient is required.")
	}

	// Set the default uid and gid so we can determine if it was
	// unset from configuration.
	tr := &TemplateResource{}

	tr.Src = config.TemplateFile
	tr.Dest = config.DestFile
	tr.Uid = uid
	tr.Gid = gid
	tr.Mode = "0644"
	tr.storeClient = config.StoreClient
	tr.funcMap = newFuncMap()
	tr.store = memkv.New()

	addFuncs(tr.funcMap, tr.store.FuncMap)

	return tr, nil
}

// setFileMode sets the FileMode.
func (t *TemplateResource) setFileMode() error {
	// check the dest file exists or not. If exists, use it's mode.Otherwise use tr.Mode=0644
	if !util.IsFileExist(t.Dest) {
		mode, err := strconv.ParseUint(t.Mode, 0, 32)
		if err != nil {
			return err
		}
		t.FileMode = os.FileMode(mode)
	} else {
		fi, err := os.Stat(t.Dest)
		if err != nil {
			return err
		}
		t.FileMode = fi.Mode()
	}

	return nil
}

// setVars sets the Vars for template resource.
func (t *TemplateResource) setVars() error {
	result, err := t.storeClient.GetValues()
	if err != nil {
		return err
	}

	t.store.Purge()

	for k, v := range result {
		t.store.Set(k, v)
	}

	return nil
}

// createStageFile stages the src configuration file by processing the src
// template and setting the desired owner, group, and mode. It also sets the
// StageFile for the template resource.
// It returns an error if any.
func (t *TemplateResource) createStageFile() error {

	tmpl, err := template.New(filepath.Base(t.Src)).Funcs(t.funcMap).ParseFiles(t.Src)
	if err != nil {
		return fmt.Errorf("Unable to process template %s, %s", t.Src, err)
	}

	if _, err := os.Stat(filepath.Dir(t.Dest)); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Dir(t.Dest), dirMode); err != nil {
			return fmt.Errorf("Create %s directory fail, error: %v ", filepath.Dir(t.Dest), err)
		}

		if err := os.Chown(filepath.Dir(t.Dest), uid, gid); err != nil {
			return fmt.Errorf("Chown %s directory fail, error: %v ", filepath.Dir(t.Dest), err)
		}
	}

	// create TempFile in Dest directory to avoid cross-filesystem issues
	temp, err := os.CreateTemp(filepath.Dir(t.Dest), "."+filepath.Base(t.Dest))
	if err != nil {
		return err
	}

	if err = tmpl.Execute(temp, nil); err != nil {
		temp.Close()
		os.Remove(temp.Name())
		return err
	}
	defer temp.Close()

	// Set the owner, group, and mode on the stage file now to make it easier to
	// compare against the destination configuration file later.
	os.Chmod(temp.Name(), t.FileMode)
	os.Chown(temp.Name(), t.Uid, t.Gid)
	t.StageFile = temp
	return nil
}

// sync compares the staged and dest config files and attempts to sync them
// if they differ. sync will run a config check command if set before
// overwriting the target config file. Finally, sync will run a reload command
// if set to have the application or service pick up the changes.
// It returns an error if any.
func (t *TemplateResource) sync() error {
	staged := t.StageFile.Name()

	defer os.Remove(staged)
	ok, err := util.IsConfigChanged(staged, t.Dest)
	if err != nil {
		return err
	}

	if ok {
		err := os.Rename(staged, t.Dest)
		if err != nil {
			if strings.Contains(err.Error(), "device or resource busy") {
				// try to open the file and write to it
				var contents []byte
				var rerr error
				contents, rerr = ioutil.ReadFile(staged)
				if rerr != nil {
					return rerr
				}
				err := ioutil.WriteFile(t.Dest, contents, t.FileMode)
				// make sure owner and group match the temp file, in case the file was created with WriteFile
				os.Chown(t.Dest, t.Uid, t.Gid)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

	}
	return nil
}

// process is a convenience function that wraps calls to the three main tasks
// required to keep local configuration files in sync. First we gather vars
// from the store, then we stage a candidate configuration file, and finally sync
// things up.
// It returns an error if any.
func (t *TemplateResource) process() error {
	if err := t.setFileMode(); err != nil {
		return err
	}

	if err := t.setVars(); err != nil {
		return err
	}

	if err := t.createStageFile(); err != nil {
		return err
	}

	return t.sync()
}