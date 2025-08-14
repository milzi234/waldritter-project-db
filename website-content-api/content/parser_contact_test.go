package content

import (
	"reflect"
	"testing"
)

func TestContactParser(t *testing.T) {
	parser := &ContactParser{}

	tests := []struct {
		name     string
		heading  string
		content  []string
		attrs    map[string]string
		expected *ComponentPageContact
	}{
		{
			name:    "full contact with social links",
			heading: "John Doe",
			content: []string{
				"Lead Organizer",
				"john.doe@waldritter.de",
				"![John Doe](/uploads/john.jpg)",
				"",
				"## Social Links",
				"- [GitHub](https://github.com/johndoe)",
				"- [LinkedIn](https://linkedin.com/in/johndoe)",
			},
			attrs: map[string]string{"id": "team-lead"},
			expected: &ComponentPageContact{
				ID:    "team-lead",
				Name:  "John Doe",
				Email: "john.doe@waldritter.de",
				Image: &Image{
					URL:     "/uploads/john.jpg",
					Caption: "John Doe",
				},
				SocialLinks: []*SocialLink{
					{
						Label:    "GitHub",
						URL:      "https://github.com/johndoe",
						External: true,
						SVGIcon:  "github",
					},
					{
						Label:    "LinkedIn",
						URL:      "https://linkedin.com/in/johndoe",
						External: true,
						SVGIcon:  "linkedin",
					},
				},
			},
		},
		{
			name:    "contact without social links",
			heading: "Jane Smith",
			content: []string{
				"jane.smith@example.com",
				"![Jane Smith](/uploads/jane.png)",
			},
			attrs: map[string]string{"id": "contact1"},
			expected: &ComponentPageContact{
				ID:    "contact1",
				Name:  "Jane Smith",
				Email: "jane.smith@example.com",
				Image: &Image{
					URL:     "/uploads/jane.png",
					Caption: "Jane Smith",
				},
				SocialLinks: []*SocialLink{},
			},
		},
		{
			name:    "contact with internal links",
			heading: "",
			content: []string{
				"Contact Person",
				"contact@site.com",
				"",
				"## Social Links",
				"[Profile](/profile)",
				"[Twitter](https://twitter.com/user)",
			},
			attrs: map[string]string{"id": "contact2"},
			expected: &ComponentPageContact{
				ID:    "contact2",
				Name:  "Contact Person",
				Email: "contact@site.com",
				SocialLinks: []*SocialLink{
					{
						Label:    "Profile",
						URL:      "/profile",
						External: false,
						SVGIcon:  "",
					},
					{
						Label:    "Twitter",
						URL:      "https://twitter.com/user",
						External: true,
						SVGIcon:  "twitter",
					},
				},
			},
		},
		{
			name:    "minimal contact",
			heading: "Bob",
			content: []string{
				"bob@example.org",
			},
			attrs: map[string]string{"id": "minimal"},
			expected: &ComponentPageContact{
				ID:          "minimal",
				Name:        "Bob",
				Email:       "bob@example.org",
				SocialLinks: []*SocialLink{},
			},
		},
		{
			name:    "detect various social icons",
			heading: "Social Test",
			content: []string{
				"test@example.com",
				"## Social Links",
				"- [Facebook Page](https://facebook.com/page)",
				"- [Instagram](https://instagram.com/user)",
				"- [YouTube Channel](https://youtube.com/channel)",
				"- [Email Me](mailto:test@example.com)",
			},
			attrs: map[string]string{"id": "social-test"},
			expected: &ComponentPageContact{
				ID:    "social-test",
				Name:  "Social Test",
				Email: "test@example.com",
				SocialLinks: []*SocialLink{
					{
						Label:    "Facebook Page",
						URL:      "https://facebook.com/page",
						External: true,
						SVGIcon:  "facebook",
					},
					{
						Label:    "Instagram",
						URL:      "https://instagram.com/user",
						External: true,
						SVGIcon:  "instagram",
					},
					{
						Label:    "YouTube Channel",
						URL:      "https://youtube.com/channel",
						External: true,
						SVGIcon:  "youtube",
					},
					{
						Label:    "Email Me",
						URL:      "mailto:test@example.com",
						External: false,
						SVGIcon:  "email",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.heading, tt.content, tt.attrs)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			contact, ok := result.(*ComponentPageContact)
			if !ok {
				t.Fatalf("Parse() returned wrong type: %T", result)
			}

			if !reflect.DeepEqual(contact, tt.expected) {
				t.Errorf("Parse() = %+v, want %+v", contact, tt.expected)
				
				// Detailed comparison for debugging
				if contact.Name != tt.expected.Name {
					t.Errorf("Name mismatch: got %q, want %q", contact.Name, tt.expected.Name)
				}
				if contact.Email != tt.expected.Email {
					t.Errorf("Email mismatch: got %q, want %q", contact.Email, tt.expected.Email)
				}
				if !reflect.DeepEqual(contact.Image, tt.expected.Image) {
					t.Errorf("Image mismatch: got %+v, want %+v", contact.Image, tt.expected.Image)
				}
				if len(contact.SocialLinks) != len(tt.expected.SocialLinks) {
					t.Errorf("SocialLinks length mismatch: got %d, want %d",
						len(contact.SocialLinks), len(tt.expected.SocialLinks))
				} else {
					for i, link := range contact.SocialLinks {
						if !reflect.DeepEqual(link, tt.expected.SocialLinks[i]) {
							t.Errorf("SocialLink[%d] mismatch: got %+v, want %+v",
								i, link, tt.expected.SocialLinks[i])
						}
					}
				}
			}
		})
	}
}

func TestContactParserCanParse(t *testing.T) {
	parser := &ContactParser{}

	tests := []struct {
		componentType string
		expected      bool
	}{
		{"contact", true},
		{"Contact", false},
		{"person", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.componentType, func(t *testing.T) {
			result := parser.CanParse(tt.componentType)
			if result != tt.expected {
				t.Errorf("CanParse(%q) = %v, want %v", tt.componentType, result, tt.expected)
			}
		})
	}
}