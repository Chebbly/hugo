// Copyright 2019 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"fmt"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/spf13/cast"
)

func createDefaultOutputFormats(allFormats output.Formats, cfg config.Provider) map[string]output.Formats {
	rssOut, rssFound := allFormats.GetByName(output.RSSFormat.Name)
	htmlOut, _ := allFormats.GetByName(output.HTMLFormat.Name)
	robotsOut, _ := allFormats.GetByName(output.RobotsTxtFormat.Name)
	sitemapOut, _ := allFormats.GetByName(output.SitemapFormat.Name)

	defaultListTypes := output.Formats{htmlOut}
	if rssFound {
		defaultListTypes = append(defaultListTypes, rssOut)
	}

	m := map[string]output.Formats{
		page.KindPage:         {htmlOut},
		page.KindHome:         defaultListTypes,
		page.KindSection:      defaultListTypes,
		page.KindTaxonomy:     defaultListTypes,
		page.KindTaxonomyTerm: defaultListTypes,
		// Below are for consistency. They are currently not used during rendering.
		kindSitemap:   {sitemapOut},
		kindRobotsTXT: {robotsOut},
		kind404:       {htmlOut},
	}

	// May be disabled
	if rssFound {
		m[kindRSS] = output.Formats{rssOut}
	}

	return m

}

func createSiteOutputFormats(allFormats output.Formats, cfg config.Provider) (map[string]output.Formats, error) {
	defaultOutputFormats := createDefaultOutputFormats(allFormats, cfg)

	if !cfg.IsSet("outputs") {
		return defaultOutputFormats, nil
	}

	outFormats := make(map[string]output.Formats)

	outputs := cfg.GetStringMap("outputs")

	if len(outputs) == 0 {
		return outFormats, nil
	}

	seen := make(map[string]bool)

	for k, v := range outputs {
		k = getKind(k)
		if k == "" {
			// Invalid kind
			continue
		}
		var formats output.Formats
		vals := cast.ToStringSlice(v)
		for _, format := range vals {
			f, found := allFormats.GetByName(format)
			if !found {
				return nil, fmt.Errorf("failed to resolve output format %q from site config", format)
			}
			formats = append(formats, f)
		}

		// This effectively prevents empty outputs entries for a given Kind.
		// We need at least one.
		if len(formats) > 0 {
			seen[k] = true
			outFormats[k] = formats
		}
	}

	// Add defaults for the entries not provided by the user.
	for k, v := range defaultOutputFormats {
		if !seen[k] {
			outFormats[k] = v
		}
	}

	return outFormats, nil

}
