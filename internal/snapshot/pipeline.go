package snapshot

import "github.com/netwatch/portwatch/internal/scanner"

// PipelineOptions configures which processing stages run.
type PipelineOptions struct {
	Normalize bool
	Dedupe    bool
	Enrich    bool
	Validate  bool
	Classify  bool
}

// DefaultPipelineOptions returns a PipelineOptions with all stages enabled.
func DefaultPipelineOptions() PipelineOptions {
	return PipelineOptions{
		Normalize: true,
		Dedupe:    true,
		Enrich:    true,
		Validate:  true,
		Classify:  true,
	}
}

// PipelineResult holds the output of a full pipeline run.
type PipelineResult struct {
	Ports      []scanner.Port
	Validation *ValidationResult
	Classified []ClassifiedPort
}

// RunPipeline executes the configured processing stages in order.
func RunPipeline(ports []scanner.Port, opts PipelineOptions) PipelineResult {
	working := make([]scanner.Port, len(ports))
	copy(working, ports)

	if opts.Normalize {
		working = Normalize(working, DefaultNormalizeOptions())
	}
	if opts.Dedupe {
		working = Dedupe(working, DefaultDedupeOptions())
	}
	if opts.Enrich {
		working = Enrich(working, DefaultEnrichOptions())
	}

	result := PipelineResult{Ports: working}

	if opts.Validate {
		vr := Validate(working, DefaultValidateOptions())
		result.Validation = &vr
	}
	if opts.Classify {
		result.Classified = Classify(working, DefaultClassifyOptions())
	}

	return result
}
