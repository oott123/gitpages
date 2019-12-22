package fileserver

import (
	"github.com/imdario/mergo"
	"github.com/oott123/gitpages/pkg/mregex"
	"net/http"
	"reflect"
)

type AccessRules []AccessConfig

type AccessConfig struct {
	// if AllowCORS is true, a Access-Control-Allow-Origin header will be set. see CORSOrigins for details.
	AllowCORS bool
	// when CORSOrigin is set, only if the request contains Origin header which matches one of the list will get ACAO header.
	// if CORSOrigin is empty list, ACAO: * will be set.
	CORSOrigins []string
	// if AllowDotFiles is true, dot files (files started with `.`) will allowed to access by user.
	AllowDotFiles bool
	// if AllowListDirectory is true, access directory which don't contains index files will result a simple list
	AllowListDirectory bool
	// if NotFoundErrorPage is defined, 404 will be sent to this file (relative to site root); else will be sent to 404.html
	NotFoundErrorPage string
	// if Match is defined, regexps should be matched before matching rules defined in current and children access config; not allowed on the top most config file
	Match *mregex.Regexp
	// Rules is children rules used for detailed access config; only allowed on the top most config file
	Rules AccessRules
	// Break will stop eval next rules
	Break bool
	// Don't use for now
	ContentType string
	// Don't use for now
	HotlinkProtection bool
	// Don't use for now
	HotlinkOrigins []string
	// Don't use for now
	ForbiddenErrorPage string
}

func (ac *AccessConfig) MatchPattern(path string) bool {
	if ac.Match == nil {
		return true
	}
	return ac.Match.MatchString(path)
}

func (ac *AccessConfig) EvaluateForPath(path string) *AccessConfig {
	if ac.Match != nil {
		log.Warnf("ignoring `Match` %s in root element; `Match` can only set on rules element", ac.Match)
	}
	if ac.Break {
		return ac
	}

	rule := ac
	for _, r := range ac.Rules {
		if len(rule.Rules) > 0 {
			log.Warnf("ignoring `Rules` %#v, `Rules` can only set on root element", rule.Rules)
		}
		if r.MatchPattern(path) {
			b := r
			rule = &b
			if rule.Break {
				break
			}
		}
	}

	return rule
}

func (ac *AccessConfig) MergeTo(dst *AccessConfig) {
	d := *dst
	err := mergo.Merge(&d, ac, mergo.WithOverride, mergo.WithTransformers(accessConfigTransformer{}))
	if err != nil {
		log.Errorf("failed to merge access config: %s", err)
	}
}

type accessConfigTransformer struct{}

func (a accessConfigTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(AccessRules{}) {
		return func(dst, src reflect.Value) error {
			return nil
		}
	}
	return nil
}

func (f *FileServer) AccessConfig(r *http.Request) *AccessConfig {
	return f.accessConfig.EvaluateForPath(r.URL.Path)
}
