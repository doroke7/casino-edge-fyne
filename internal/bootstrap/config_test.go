package bootstrap

import (
	"os"
	"path/filepath"
	"testing"
)

func writeYAML(t *testing.T, sDirectory, sName, sContent string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(sDirectory, sName), []byte(sContent), 0o644); err != nil {
		t.Fatal(err)
	}
}

const sValidCamera = `
width: 1280
height: 720
framerate: 30
record_directory: '~/Movies'
record_bitrate: '4M'
preview_width: 640
preview_height: 360
preview_framerate: 15
snapshot_directory: './runtime/snapshots'
snapshot_interval: 600
`

func TestLoadUsesFileNameAsNamespace(t *testing.T) {
	sDirectory := t.TempDir()
	writeYAML(t, sDirectory, "camera.yaml", sValidCamera)
	writeYAML(t, sDirectory, "default.yaml", "debug: true\n")

	oConfig, err := Load(sDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if oConfig.CAMERA.WIDTH != 1280 || oConfig.CAMERA.RECORD_BITRATE != "4M" || oConfig.CAMERA.SNAPSHOT_INTERVAL != 600 {
		t.Fatalf("camera.yaml not mapped under camera.*: %+v", oConfig.CAMERA)
	}
	if !oConfig.DEFAULT.DEBUG {
		t.Fatal("default.yaml not mapped under default.*")
	}
}

func TestLoadEnvOverridesNestedField(t *testing.T) {
	sDirectory := t.TempDir()
	writeYAML(t, sDirectory, "camera.yaml", sValidCamera)

	t.Setenv("CAMERA_SNAPSHOT_INTERVAL", "5")
	t.Setenv("CAMERA_RECORD_DIRECTORY", "/tmp/somewhere")

	oConfig, err := Load(sDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if oConfig.CAMERA.SNAPSHOT_INTERVAL != 5 || oConfig.CAMERA.RECORD_DIRECTORY != "/tmp/somewhere" {
		t.Fatalf("env did not override yaml: %+v", oConfig.CAMERA)
	}
	if oConfig.CAMERA.WIDTH != 1280 {
		t.Fatal("non-overridden field must keep yaml value")
	}
}

func TestValidateRejectsMissingCameraConfig(t *testing.T) {
	oConfig, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := oConfig.Validate(); err == nil {
		t.Fatal("empty config must fail validation")
	}
}

func TestValidateAcceptsCompleteConfig(t *testing.T) {
	sDirectory := t.TempDir()
	writeYAML(t, sDirectory, "camera.yaml", sValidCamera)

	oConfig, err := Load(sDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if err := oConfig.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestLoadRejectsBrokenYAML(t *testing.T) {
	sDirectory := t.TempDir()
	writeYAML(t, sDirectory, "camera.yaml", "width: [1280\n")

	if _, err := Load(sDirectory); err == nil {
		t.Fatal("broken yaml must be rejected")
	}
}
