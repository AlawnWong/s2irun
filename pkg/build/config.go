package build

import (
	"errors"
	"fmt"

	"github.com/kubesphere/s2irun/pkg/api"
	"github.com/kubesphere/s2irun/pkg/api/constants"
	"github.com/kubesphere/s2irun/pkg/docker"
	"github.com/kubesphere/s2irun/pkg/scm/git"
)

// GenerateConfigFromLabels generates the S2I Config struct from the Docker
// image labels.
func GenerateConfigFromLabels(config *api.Config, metadata *docker.PullResult) error {
	if config == nil {
		return errors.New("config must be provided to GenerateConfigFromLabels")
	}
	if metadata == nil {
		return errors.New("image metadata must be provided to GenerateConfigFromLabels")
	}

	labels := metadata.Image.Config.Labels

	if builderVersion, ok := labels[constants.BuilderVersionLabel]; ok {
		config.BuilderImageVersion = builderVersion
		config.BuilderBaseImageVersion = labels[constants.BuilderBaseVersionLabel]
	}

	config.ScriptsURL = labels[constants.ScriptsURLLabel]
	if len(config.ScriptsURL) == 0 {
		// FIXME: Backward compatibility
		config.ScriptsURL = labels[constants.DeprecatedScriptsURLLabel]
	}

	config.Description = labels[constants.KubernetesDescriptionLabel]
	config.DisplayName = labels[constants.KubernetesDisplayNameLabel]

	if builder, ok := labels[constants.BuildImageLabel]; ok {
		config.BuilderImage = builder
	} else {
		return fmt.Errorf("required label %q not found in image", constants.BuildImageLabelLog)
	}

	if repo, ok := labels[constants.BuildSourceLocationLabel]; ok {
		source, err := git.Parse(repo, false)
		if err != nil {
			return fmt.Errorf("couldn't parse label %q value %s: %v", constants.BuildSourceLocationLabelLog, repo, err)
		}
		config.Source = source
	} else {
		return fmt.Errorf("required label %q not found in image", constants.BuildSourceLocationLabelLog)
	}

	config.ContextDir = labels[constants.BuildSourceContextDirLabel]
	config.Source.URL.Fragment = labels[constants.BuildCommitRefLabel]

	return nil
}
