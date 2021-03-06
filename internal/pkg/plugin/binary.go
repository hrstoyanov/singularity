// Copyright (c) 2018, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"plugin"

	"github.com/sylabs/sif/pkg/sif"
	"github.com/sylabs/singularity/internal/pkg/util/fs"
	pluginapi "github.com/sylabs/singularity/pkg/plugin"
)

const (
	// DirRoot is the root directory for the plugin installation, typically
	// located within LIBEXECDIR.
	DirRoot = "plugin"
	// NameImage is the name of the SIF image of the plugin
	NameImage = "plugin.sif"
	// NameBinary is the name of the plugin object
	NameBinary = "object.so"
	// NameConfig is the name of the plugin's config file
	NameConfig = "config.yaml"
)

// Meta is an internal representation of a plugin binary and all of its
// artifacts. This represents the on-disk location of the SIF, shared library,
// config file, etc... This struct is written as JSON into the DirRoot directory.
type Meta struct {
	// Name is the name of the plugin
	Name string
	// Path is a path, derived from its Name, which the plugins
	// artifacts (config, SIF, .so, etc...) are located
	Path string
	// Enabled reports whether or not the plugin should be loaded
	Enabled bool

	fimg   *sif.FileImage // Plugin SIF object
	binary *plugin.Plugin // Plugin binary object
	cfg    *os.File       // Plugin YAML config file

	file *os.File // Pointer to Meta file on disk, for Read/Write access
}

// LoadFromJSON loads a Meta type from an io.Reader containing JSON. A plugin Meta
// object created in this form is read-only.
func LoadFromJSON(r io.Reader) (*Meta, error) {
	m := &Meta{}

	if err := json.NewDecoder(r).Decode(m); err != nil {
		return nil, fmt.Errorf("while decoding Meta JSON file: %s", err)
	}

	m.cfg, _ = m.Config()

	return m, nil
}

// Config returns the plugin configuration file opened as an os.File object
func (m *Meta) Config() (*os.File, error) {
	if !fs.IsFile(m.configName()) {
		return nil, nil
	}

	return os.Open(m.configName())
}

// NewFromImage returns a new meta object which hasn't yet been installed from
// a pointer to an on disk SIF. It will:
//     1. Check that the SIF is a valid plugin
//     2. Open the Manifest to retrieve name and calculate the path
//     3. Copy the SIF into the plugin path
//     4. Extract the binary object into the path
//     5. Generate a default config file in the path
//     6. Write the Meta struct onto disk in DirRoot
func NewFromImage(fimg *sif.FileImage, libexecdir string) (*Meta, error) {
	if !isPluginFile(fimg) {
		return nil, fmt.Errorf("while opening sif file: not a valid plugin")
	}

	manifest := getManifest(fimg)
	abspath, err := filepath.Abs(filepath.Join(libexecdir, pathFromManifest(manifest)))
	if err != nil {
		return nil, fmt.Errorf("while getting absolute path to plugin installation: %s", err)
	}

	m := &Meta{
		Name:    manifest.Name,
		Path:    abspath,
		Enabled: true,

		fimg: fimg,
	}

	m.installTo(libexecdir)
	return m, nil
}

// installTo installs the plugin represented by m into libexecdir. This should
// normally only be called in NewFromImage
func (m *Meta) installTo(libexecdir string) {

}

//
// Misc helper functions
//

// pathFromManifest returns a path which will exist inside of DirRoot and
// is derived from Manifest.Name
func pathFromManifest(pluginapi.Manifest) string {
	return ""
}

// metaFileFromName returns the name of the Meta file from the plugin name, which
// is a unique string generated by hashing n
func metaFileFromName(n string) string {
	return ""
}

// copyFile copies a file from src -> dst
func copyFile(src, dst string) error {
	// copycmd := exec.Command("cp", src, dst)
	return nil
}

//
// Path name helper methods on (m *Meta)
//

func (m *Meta) imageName() string {
	return filepath.Join(m.Path, NameImage)
}

func (m *Meta) binaryName() string {
	return filepath.Join(m.Path, NameBinary)
}

func (m *Meta) configName() string {
	return filepath.Join(m.Path, NameConfig)
}

//
// Helper functions for fimg *sif.FileImage
//

// isPluginFile checks if the sif.FileImage contains the sections which
// make up a valid plugin. A plugin sif file should have the following
// format:
//
// DESCR[0]: Sifplugin
//   - Datatype: sif.DataPartition
//   - Fstype:   sif.FsRaw
//   - Parttype: sif.PartData
// DESCR[1]: Sifmanifest
//   - Datatype: sif.DataGenericJSON
func isPluginFile(fimg *sif.FileImage) bool {
	return false
}

// getManifest will extract the Manifest data from the input FileImage
func getManifest(fimg *sif.FileImage) pluginapi.Manifest {
	return pluginapi.Manifest{}
}
