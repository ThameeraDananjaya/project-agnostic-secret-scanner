package gitleaks

import (
	"bufio"
	"errors"
	"os"
	"regexp/syntax"
	"strings"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine"
)

const RequiredRuleCount = 224

type RuleSpanProof struct {
	RuleCount int
	Maximum   int64
}

// ProveRuleSpans parses every single-line TOML rule regex using the exact Go
// standard-library regex grammar selected for the pinned detector build.
func ProveRuleSpans(path, digest string) (RuleSpanProof, error) {
	if err := engine.VerifyRegularFile(path, digest); err != nil {
		return RuleSpanProof{}, err
	}
	f, err := os.Open(path)
	if err != nil {
		return RuleSpanProof{}, errors.New("rule pack cannot be opened")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 4096), 1<<20)
	proof := RuleSpanProof{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "regex = '''") {
			continue
		}
		if !strings.HasSuffix(line, "'''") {
			return RuleSpanProof{}, errors.New("multiline or malformed rule regex")
		}
		pattern := strings.TrimSuffix(strings.TrimPrefix(line, "regex = '''"), "'''")
		re, err := syntax.Parse(pattern, syntax.Perl)
		if err != nil {
			return RuleSpanProof{}, errors.New("rule regex is unsupported")
		}
		span := regexpByteSpan(re)
		if span.max < 0 {
			return RuleSpanProof{}, errors.New("rule regex has unbounded byte span")
		}
		proof.RuleCount++
		if span.max > proof.Maximum {
			proof.Maximum = span.max
		}
	}
	if scanner.Err() != nil || proof.RuleCount != RequiredRuleCount || proof.Maximum != 4020 {
		return RuleSpanProof{}, errors.New("rule span proof binding mismatch")
	}
	return proof, nil
}

type byteSpan struct{ min, max int64 }

func addSpan(a, b byteSpan) byteSpan {
	maximum := int64(-1)
	if a.max >= 0 && b.max >= 0 {
		maximum = a.max + b.max
	}
	return byteSpan{min: a.min + b.min, max: maximum}
}

func regexpByteSpan(re *syntax.Regexp) byteSpan {
	switch re.Op {
	case syntax.OpNoMatch, syntax.OpEmptyMatch, syntax.OpBeginLine, syntax.OpEndLine,
		syntax.OpBeginText, syntax.OpEndText, syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		return byteSpan{}
	case syntax.OpLiteral:
		var size int64
		for _, r := range re.Rune {
			switch {
			case r <= 0x7f:
				size++
			case r <= 0x7ff:
				size += 2
			case r <= 0xffff:
				size += 3
			default:
				size += 4
			}
		}
		return byteSpan{min: size, max: size}
	case syntax.OpCharClass, syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		return byteSpan{min: 1, max: 4}
	case syntax.OpCapture:
		return regexpByteSpan(re.Sub[0])
	case syntax.OpConcat:
		span := byteSpan{}
		for _, sub := range re.Sub {
			span = addSpan(span, regexpByteSpan(sub))
		}
		return span
	case syntax.OpAlternate:
		span := byteSpan{min: (1 << 62) - 1, max: 0}
		for _, sub := range re.Sub {
			candidate := regexpByteSpan(sub)
			if candidate.min < span.min {
				span.min = candidate.min
			}
			if candidate.max < 0 || span.max < 0 {
				span.max = -1
			} else if candidate.max > span.max {
				span.max = candidate.max
			}
		}
		return span
	case syntax.OpQuest:
		span := regexpByteSpan(re.Sub[0])
		span.min = 0
		return span
	case syntax.OpStar:
		span := regexpByteSpan(re.Sub[0])
		if span.max == 0 {
			return byteSpan{}
		}
		return byteSpan{max: -1}
	case syntax.OpPlus:
		span := regexpByteSpan(re.Sub[0])
		if span.max == 0 {
			return byteSpan{}
		}
		return byteSpan{min: span.min, max: -1}
	case syntax.OpRepeat:
		span := regexpByteSpan(re.Sub[0])
		if re.Max < 0 && span.max > 0 {
			return byteSpan{min: span.min * int64(re.Min), max: -1}
		}
		return byteSpan{min: span.min * int64(re.Min), max: span.max * int64(re.Max)}
	default:
		return byteSpan{max: -1}
	}
}
