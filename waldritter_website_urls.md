# Waldritter Website URL Structure - Refreshed Site Plan

Based on analysis of the current waldritter.de website, here's a comprehensive list of URLs for the refreshed site with checkboxes to track implementation:

## Core Pages

### Main Navigation
- [ ] `/` - Homepage (Start)
  - Hero section with mission statement
  - Featured projects/news
  - Newsletter signup
  - Quick links to key areas
  
- [ ] `/leitbild` - Guiding Principles (Leitbild)
  - Organization values and mission
  - Educational philosophy
  - Diversity and inclusion statement
  
- [ ] `/ueber-uns` - About Us
  - History of Waldritter
  - Team/Board members
  - Organizational structure
  - Annual reports
  
- [ ] `/kontakt` - Contact
  - Contact form
  - Office hours
  - Address and phone
  - Map/directions

## Local Groups (Ortsgruppen)

- [ ] `/ortsgruppen` - Local Groups Overview
  - Interactive map of locations
  - List of all groups
  
### Individual Group Pages
- [ ] `/ortsgruppen/berlin` - Berlin Group
- [ ] `/ortsgruppen/freudenberg` - Freudenberg Group
- [ ] `/ortsgruppen/giessen` - Gießen Group
- [ ] `/ortsgruppen/herten` - Herten Group
- [ ] `/ortsgruppen/lahn-dill` - Lahn-Dill Group
- [ ] `/ortsgruppen/potsdam` - Potsdam Group
- [ ] `/ortsgruppen/siegen` - Siegen Group
- [ ] `/ortsgruppen/sued-west` - Süd-West Group
- [ ] `/ortsgruppen/westerwald` - Westerwald Group
- [ ] `/ortsgruppen/witten` - Witten Group

## Educational Offerings

### School Programs
- [ ] `/schulangebote` - School Offerings Overview
- [ ] `/schulangebote/grundschule` - Elementary School Programs
- [ ] `/schulangebote/sekundarstufe-1` - Lower Secondary Programs
- [ ] `/schulangebote/sekundarstufe-2` - Upper Secondary Programs
- [ ] `/schulangebote/berufskolleg` - Vocational College Programs
- [ ] `/schulangebote/waldprogramme` - Forest Experience Programs
- [ ] `/schulangebote/buchung` - Booking Form

### Youth Programs
- [ ] `/jugendprogramme` - Youth Programs Overview
- [ ] `/jugendprogramme/6-10-jahre` - Programs for Ages 6-10
- [ ] `/jugendprogramme/11-14-jahre` - Programs for Ages 11-14
- [ ] `/jugendprogramme/15-16-jahre` - Programs for Ages 15-16
- [ ] `/jugendprogramme/ferienfreizeiten` - Holiday Camps

## Projects & Initiatives

- [ ] `/projekte` - Projects Overview
- [ ] `/projekte/larp-fuer-demokratie` - LARP for Democracy Project
- [ ] `/projekte/tagungshaus-herten` - Herten Meeting House
- [ ] `/projekte/modellprojekte` - Model Projects
- [ ] `/projekte/geförderte-projekte` - Funded Projects

## Resources & Information

### For Participants
- [ ] `/teilnehmer` - Participant Information
- [ ] `/teilnehmer/erste-schritte` - Getting Started
- [ ] `/teilnehmer/faq` - Frequently Asked Questions
- [ ] `/teilnehmer/elterninfo` - Parent Information
- [ ] `/teilnehmer/sicherheit` - Safety Information

### For Educators
- [ ] `/paedagogen` - Educator Resources
- [ ] `/paedagogen/methoden` - Educational Methods
- [ ] `/paedagogen/materialien` - Teaching Materials
- [ ] `/paedagogen/fortbildungen` - Professional Development

### Publications & Media
- [ ] `/publikationen` - Publications
- [ ] `/publikationen/newsletter` - Newsletter Archive
- [ ] `/publikationen/jahresberichte` - Annual Reports
- [ ] `/publikationen/forschung` - Research Papers
- [ ] `/medien/galerie` - Photo/Video Gallery
- [ ] `/medien/presse` - Press Kit

## Support & Involvement

- [ ] `/mitmachen` - Get Involved
- [ ] `/mitmachen/mitgliedschaft` - Membership
- [ ] `/mitmachen/ehrenamt` - Volunteer Opportunities
- [ ] `/mitmachen/praktikum` - Internships
- [ ] `/mitmachen/jobs` - Job Openings

- [ ] `/spenden` - Donations
- [ ] `/spenden/einzelspende` - One-time Donation
- [ ] `/spenden/foerdermitglied` - Supporting Member
- [ ] `/spenden/sponsoring` - Corporate Sponsorship

## Legal & Administrative

- [ ] `/impressum` - Legal Notice
- [ ] `/datenschutz` - Privacy Policy
- [ ] `/agb` - Terms and Conditions
- [ ] `/barrierefreiheit` - Accessibility Statement
- [ ] `/sitemap` - Site Map

### Special Content
- [ ] `/partner` - Partners & Sponsors
- [ ] `/netzwerk` - Network & Collaborations

---

## Implementation Plan with Markdown-Based Content API

### System Architecture
- **Content Management**: Markdown files in `website-content-api/obsidian-vault/pages/`
- **Routing**: Automatic - Vue router queries GraphQL API for all pages
- **Dynamic Content**: Rails project database for events, activities, galleries
- **Navigation**: Defined in `obsidian-vault/navigation/main-menu.md`

### Implementation Process

#### 1. **Content Creation** (Using waldritter-content-writer agent)
Each page will be a markdown file with:
- `documentId`: Unique identifier
- `url`: The route path (automatically creates Vue route)
- `components`: Hero, markdown, contact, projectSearch, highlightReel
- German content based on current waldritter.de
- Images downloaded from waldritter.de and stored in `obsidian-vault/assets/images/uploads/`

#### 2. **Priority Pages** (Phase 1):
- [x] Homepage - Update existing with Waldritter content
- [x] `/leitbild` - Core values and mission
- [x] `/kontakt` - Contact information
- [x] `/ueber-uns` - About Waldritter
- [x] `/ortsgruppen` - Local groups overview
- [x] `/schulangebote` - School offerings

#### 3. **Secondary Pages** (Phase 2):
- [ ] Individual Ortsgruppen pages (10 local groups)
- [ ] Educational program detail pages
- [ ] `/teilnehmer/*` - Participant information
- [ ] `/mitmachen/*` - Involvement opportunities

#### 4. **Dynamic Content Integration**:
- **Events/Termine**: Display via projectSearch component (type: event)
- **Local Activities**: Projects tagged by location
- **Galleries**: Project images with gallery tag
- **Educational Programs**: Projects with education category

#### 5. **Navigation Updates**:
- [x] Update `main-menu.md` with German labels
- [x] Add hierarchical menu structure
- [x] Update `footer.md` with links and social media

#### 6. **Image Migration Process**:
- Download images from waldritter.de pages as we create content
- Store in `website-content-api/obsidian-vault/assets/images/uploads/`
- Convert to WebP format for optimal performance
- Maintain proper licensing and attribution
- Reference in markdown as `/uploads/filename.webp`

### Technical Benefits
- ✅ No router configuration needed (automatic from GraphQL)
- ✅ Hot reload on content changes
- ✅ Version control for all content (Git)
- ✅ Edit with Obsidian or any text editor
- ✅ 5-10x faster than Strapi
- ✅ No database required for static content

### Content Component Types Available
- **hero**: Hero sections with image and CTA
- **markdown**: Rich text content
- **contact**: Contact forms and information
- **projectSearch**: Dynamic project/event listings
- **highlightReel**: Featured content carousel

This approach leverages the markdown-based CMS while maintaining full compatibility with the existing Vue frontend and Rails project database.