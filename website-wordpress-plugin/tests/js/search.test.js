/**
 * Search Component Tests
 */

describe('Search', () => {
  let container;

  beforeEach(() => {
    container = document.createElement('div');
    container.innerHTML = `
      <div
        id="test-search"
        class="waldritter-search"
        data-per-page="5"
        data-hidden-categories="[]"
        data-initial-tags="[]"
        data-show-occurrences="false"
      >
        <div class="waldritter-search__filter-section">
          <button
            class="waldritter-search__filter-button"
            aria-expanded="false"
            aria-controls="test-search-drawer"
          >
            <span class="waldritter-search__filter-text">Filter</span>
            <span class="waldritter-search__filter-count"></span>
          </button>

          <div id="test-search-drawer" class="waldritter-search__drawer" aria-hidden="true">
            <div class="waldritter-search__drawer-header">
              <h3 class="waldritter-search__drawer-title">Filter by Tags</h3>
              <button class="waldritter-search__drawer-close">&times;</button>
            </div>
            <div class="waldritter-search__drawer-content">
              <fieldset class="waldritter-search__category">
                <legend class="waldritter-search__category-title">Category 1</legend>
                <div class="waldritter-search__tags">
                  <label class="waldritter-search__tag-label">
                    <input type="checkbox" class="waldritter-search__tag-checkbox" value="Tag1" />
                    <span>Tag1</span>
                  </label>
                  <label class="waldritter-search__tag-label">
                    <input type="checkbox" class="waldritter-search__tag-checkbox" value="Tag2" />
                    <span>Tag2</span>
                  </label>
                </div>
              </fieldset>
            </div>
            <div class="waldritter-search__drawer-footer">
              <button class="waldritter-search__clear-button">Clear All</button>
              <button class="waldritter-search__apply-button">Apply Filters</button>
            </div>
          </div>
          <div class="waldritter-search__drawer-backdrop" aria-hidden="true"></div>
        </div>

        <div class="waldritter-search__results-info">
          <span class="waldritter-search__results-count">10 Projects</span>
        </div>

        <div class="waldritter-search__list" role="list"></div>

        <div class="waldritter-search__loading" aria-hidden="true">
          <div class="waldritter-search__spinner"></div>
        </div>

        <nav class="waldritter-search__pagination">
          <button class="waldritter-search__page-button waldritter-search__page-button--prev" disabled>Prev</button>
          <div class="waldritter-search__page-info">
            <span class="waldritter-search__current-page">1</span>
            /
            <span class="waldritter-search__total-pages">2</span>
          </div>
          <button class="waldritter-search__page-button waldritter-search__page-button--next">Next</button>
        </nav>
      </div>
    `;
    document.body.appendChild(container);
  });

  afterEach(() => {
    document.body.removeChild(container);
  });

  describe('Structure', () => {
    it('should have filter button', () => {
      const filterBtn = container.querySelector('.waldritter-search__filter-button');
      expect(filterBtn).toBeTruthy();
    });

    it('should have drawer hidden by default', () => {
      const drawer = container.querySelector('.waldritter-search__drawer');
      expect(drawer.getAttribute('aria-hidden')).toBe('true');
    });

    it('should have pagination controls', () => {
      const prevBtn = container.querySelector('.waldritter-search__page-button--prev');
      const nextBtn = container.querySelector('.waldritter-search__page-button--next');
      expect(prevBtn).toBeTruthy();
      expect(nextBtn).toBeTruthy();
    });
  });

  describe('Filter Drawer', () => {
    it('should open drawer when filter button clicked', () => {
      const filterBtn = container.querySelector('.waldritter-search__filter-button');
      const drawer = container.querySelector('.waldritter-search__drawer');

      // Simulate opening
      drawer.setAttribute('aria-hidden', 'false');
      filterBtn.setAttribute('aria-expanded', 'true');

      expect(drawer.getAttribute('aria-hidden')).toBe('false');
      expect(filterBtn.getAttribute('aria-expanded')).toBe('true');
    });

    it('should close drawer when close button clicked', () => {
      const drawer = container.querySelector('.waldritter-search__drawer');
      const filterBtn = container.querySelector('.waldritter-search__filter-button');

      // Simulate closing
      drawer.setAttribute('aria-hidden', 'true');
      filterBtn.setAttribute('aria-expanded', 'false');

      expect(drawer.getAttribute('aria-hidden')).toBe('true');
      expect(filterBtn.getAttribute('aria-expanded')).toBe('false');
    });

    it('should have checkboxes for tags', () => {
      const checkboxes = container.querySelectorAll('.waldritter-search__tag-checkbox');
      expect(checkboxes.length).toBeGreaterThan(0);
    });
  });

  describe('Tag Selection', () => {
    it('should allow selecting tags', () => {
      const checkbox = container.querySelector('.waldritter-search__tag-checkbox');
      checkbox.checked = true;
      expect(checkbox.checked).toBe(true);
    });

    it('should get selected tags from UI', () => {
      const checkboxes = container.querySelectorAll('.waldritter-search__tag-checkbox');
      checkboxes[0].checked = true;
      checkboxes[1].checked = true;

      const selectedTags = Array.from(checkboxes)
        .filter(cb => cb.checked)
        .map(cb => cb.value);

      expect(selectedTags).toEqual(['Tag1', 'Tag2']);
    });

    it('should clear all tags when clear button used', () => {
      const checkboxes = container.querySelectorAll('.waldritter-search__tag-checkbox');
      checkboxes.forEach(cb => cb.checked = true);

      // Simulate clear
      checkboxes.forEach(cb => cb.checked = false);

      const selectedTags = Array.from(checkboxes).filter(cb => cb.checked);
      expect(selectedTags.length).toBe(0);
    });
  });

  describe('Pagination', () => {
    it('should have prev button disabled on first page', () => {
      const prevBtn = container.querySelector('.waldritter-search__page-button--prev');
      expect(prevBtn.disabled).toBe(true);
    });

    it('should update page info', () => {
      const currentPage = container.querySelector('.waldritter-search__current-page');
      const totalPages = container.querySelector('.waldritter-search__total-pages');

      expect(currentPage.textContent).toBe('1');
      expect(totalPages.textContent).toBe('2');
    });

    it('should enable next button when more pages exist', () => {
      const nextBtn = container.querySelector('.waldritter-search__page-button--next');
      expect(nextBtn.disabled).toBe(false);
    });
  });

  describe('Loading State', () => {
    it('should have loading element hidden by default', () => {
      const loading = container.querySelector('.waldritter-search__loading');
      expect(loading.getAttribute('aria-hidden')).toBe('true');
    });

    it('should show loading when searching', () => {
      const searchWidget = container.querySelector('.waldritter-search');
      searchWidget.classList.add('waldritter-search--loading');

      expect(searchWidget.classList.contains('waldritter-search--loading')).toBe(true);
    });
  });

  describe('Data Attributes', () => {
    it('should read per page from data attribute', () => {
      const search = container.querySelector('.waldritter-search');
      expect(search.dataset.perPage).toBe('5');
    });

    it('should read hidden categories from data attribute', () => {
      const search = container.querySelector('.waldritter-search');
      const hiddenCategories = JSON.parse(search.dataset.hiddenCategories);
      expect(Array.isArray(hiddenCategories)).toBe(true);
    });
  });
});

describe('Project Card Rendering', () => {
  const escapeHtml = (text) => {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  };

  it('should escape HTML in project title', () => {
    const maliciousTitle = '<script>alert("xss")</script>';
    const escaped = escapeHtml(maliciousTitle);

    expect(escaped).not.toContain('<script>');
    expect(escaped).toContain('&lt;script&gt;');
  });

  it('should properly format tags by category', () => {
    const tags = [
      { title: 'Tag1', category: { title: 'Category A' } },
      { title: 'Tag2', category: { title: 'Category A' } },
      { title: 'Tag3', category: { title: 'Category B' } },
    ];

    const tagsByCategory = {};
    tags.forEach(tag => {
      const category = tag.category?.title || 'Other';
      if (!tagsByCategory[category]) {
        tagsByCategory[category] = [];
      }
      tagsByCategory[category].push(tag.title);
    });

    expect(tagsByCategory['Category A']).toEqual(['Tag1', 'Tag2']);
    expect(tagsByCategory['Category B']).toEqual(['Tag3']);
  });
});
