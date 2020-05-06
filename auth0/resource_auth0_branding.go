package auth0

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gopkg.in/auth0.v4/management"
)

func newBranding() *schema.Resource {
	return &schema.Resource{
		Create: createBranding,
		Read:   readBranding,
		Update: updateBranding,
		Delete: deleteBranding,

		Schema: map[string]*schema.Schema{
			"colors": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"page_background": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"colors.0.page_background_gradient"},
						},
						"page_background_gradient": {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: []string{"colors.0.page_background"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"start": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"end": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"angle_deg": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"favicon_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logo_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"font": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func createBranding(d *schema.ResourceData, m interface{}) error {
	d.SetId(resource.UniqueId())
	return updateBranding(d, m)
}

func readBranding(d *schema.ResourceData, m interface{}) error {
	api := m.(*management.Management)

	b, err := api.Branding.Read()
	if err != nil {
		if mErr, ok := err.(management.Error); ok {
			if mErr.Status() == http.StatusNotFound {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("colors", flattenBrandingColors(b.Colors))
	d.Set("favicon_url", b.FaviconURL)
	d.Set("logo_url", b.LogoURL)
	d.Set("font", flattenBrandingFont(b.Font))

	return nil
}

func updateBranding(d *schema.ResourceData, m interface{}) error {
	b := buildBranding(d)
	api := m.(*management.Management)
	err := api.Branding.Update(b)
	if err != nil {
		return err
	}
	return readBranding(d, m)
}

func deleteBranding(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func buildBranding(d *schema.ResourceData) *management.Branding {
	b := &management.Branding{
		Colors:     expandBrandingColors(d),
		FaviconURL: String(d, "favicon_url"),
		LogoURL:    String(d, "logo_url"),
		Font:       expandBrandingFont(d),
	}

	return b
}