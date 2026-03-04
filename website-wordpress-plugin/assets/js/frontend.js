/**
 * Waldritter Project DB - Frontend JavaScript
 *
 * Handles carousel navigation, search filtering, and pagination.
 */

(function() {
  'use strict';

  /**
   * Markdown renderer using marked.js
   * Configured with safe defaults to prevent XSS
   */
  const markdown = {
    init() {
      if (typeof window.marked !== 'undefined') {
        window.marked.setOptions({
          breaks: true,      // Convert \n to <br>
          gfm: true,         // GitHub Flavored Markdown
          headerIds: false,  // Don't generate header IDs
          mangle: false,     // Don't mangle email addresses
        });
      }
    },

    /**
     * Render markdown to HTML with XSS sanitization
     * @param {string} text - Markdown text
     * @returns {string} - Sanitized HTML
     */
    render(text) {
      if (!text) return '';
      if (typeof window.marked === 'undefined') {
        // Fallback: escape HTML and convert newlines to <br>
        return this.escapeHtml(text).replace(/\n/g, '<br>');
      }
      // Parse markdown and sanitize
      const html = window.marked.parse(text);
      return this.sanitize(html);
    },

    /**
     * Escape HTML special characters
     * @param {string} text
     * @returns {string}
     */
    escapeHtml(text) {
      const div = document.createElement('div');
      div.textContent = text;
      return div.innerHTML;
    },

    /**
     * Basic sanitization - allow only safe HTML tags
     * @param {string} html
     * @returns {string}
     */
    sanitize(html) {
      const allowedTags = ['p', 'br', 'strong', 'b', 'em', 'i', 'a', 'ul', 'ol', 'li', 'code', 'pre', 'blockquote', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6'];
      const allowedAttributes = { 'a': ['href', 'target', 'rel'] };

      const div = document.createElement('div');
      div.innerHTML = html;

      // Walk through all elements and remove disallowed ones
      const walker = document.createTreeWalker(div, NodeFilter.SHOW_ELEMENT, null, false);
      const nodesToRemove = [];

      while (walker.nextNode()) {
        const node = walker.currentNode;
        const tagName = node.tagName.toLowerCase();

        if (!allowedTags.includes(tagName)) {
          // Replace disallowed tag with its text content
          nodesToRemove.push(node);
        } else {
          // Remove disallowed attributes
          const allowed = allowedAttributes[tagName] || [];
          Array.from(node.attributes).forEach(attr => {
            if (!allowed.includes(attr.name)) {
              node.removeAttribute(attr.name);
            }
          });
          // Ensure links are safe
          if (tagName === 'a') {
            const href = node.getAttribute('href') || '';
            if (href.startsWith('javascript:') || href.startsWith('data:')) {
              node.removeAttribute('href');
            }
            node.setAttribute('target', '_blank');
            node.setAttribute('rel', 'noopener noreferrer');
          }
        }
      }

      // Remove disallowed nodes (replace with text content)
      nodesToRemove.forEach(node => {
        const text = document.createTextNode(node.textContent);
        node.parentNode.replaceChild(text, node);
      });

      return div.innerHTML;
    }
  };

  /**
   * Date formatting helper
   */
  const dateHelper = {
    /**
     * Format a date range for display
     * @param {string} startDate - ISO8601 start datetime
     * @param {string|null} endDate - ISO8601 end datetime
     * @returns {string} Formatted date range
     */
    formatRange(startDate, endDate = null) {
      if (!startDate) return '';

      try {
        const start = new Date(startDate);
        const startDate = this.formatDate(start);
        const startTime = this.formatTime(start);

        if (!endDate) {
          return startTime !== '00:00' ? `${startDate} ${startTime}` : startDate;
        }

        const end = new Date(endDate);
        const endDate = this.formatDate(end);
        const endTime = this.formatTime(end);

        // Same day
        if (startDate === endDate) {
          if (startTime !== '00:00' || endTime !== '00:00') {
            return `${startDate} ${startTime} - ${endTime}`;
          }
          return startDate;
        }

        // Different days
        if (startTime !== '00:00' || endTime !== '00:00') {
          return `${startDate} ${startTime} - ${endDate} ${endTime}`;
        }
        return `${startDate} - ${endDate}`;
      } catch (e) {
        return '';
      }
    },

    /**
     * Format date as DD.MM.YYYY
     */
    formatDate(date) {
      const day = String(date.getDate()).padStart(2, '0');
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const year = date.getFullYear();
      return `${day}.${month}.${year}`;
    },

    /**
     * Format time as HH:MM
     */
    formatTime(date) {
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      return `${hours}:${minutes}`;
    },

    /**
     * Get relative date string (Today, Tomorrow, In X days)
     */
    getRelative(dateTime) {
      if (!dateTime) return null;

      try {
        const date = new Date(dateTime);
        const now = new Date();

        // Reset time part for day comparison
        const dateDay = new Date(date.getFullYear(), date.getMonth(), date.getDate());
        const nowDay = new Date(now.getFullYear(), now.getMonth(), now.getDate());

        const diffTime = dateDay - nowDay;
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

        if (diffDays < 0) return null; // Past
        if (diffDays === 0) return 'Heute';
        if (diffDays === 1) return 'Morgen';
        if (diffDays < 7) return `In ${diffDays} Tagen`;
        if (diffDays < 30) {
          const weeks = Math.floor(diffDays / 7);
          return weeks === 1 ? 'In 1 Woche' : `In ${weeks} Wochen`;
        }
        return null;
      } catch (e) {
        return null;
      }
    },

    /**
     * Check if datetime is in the future
     */
    isFuture(dateTime) {
      if (!dateTime) return false;
      try {
        return new Date(dateTime) > new Date();
      } catch (e) {
        return false;
      }
    },

    /**
     * Filter and get upcoming occurrences
     */
    getUpcoming(occurrences, limit = 3) {
      if (!occurrences || !Array.isArray(occurrences)) return [];

      const now = new Date();
      return occurrences
        .filter(occ => {
          if (!occ.startDate) return false;
          try {
            return new Date(occ.startDate) >= now;
          } catch (e) {
            return false;
          }
        })
        .slice(0, limit);
    }
  };

  /**
   * API client for REST endpoints
   */
  const api = {
    baseUrl: window.waldritterProjectDB?.restUrl || '/wp-json/waldritter/v1/',
    nonce: window.waldritterProjectDB?.nonce || '',

    async fetch(endpoint, options = {}) {
      // Handle URL construction when baseUrl already contains query string (non-pretty permalinks)
      let url;
      if (this.baseUrl.includes('?')) {
        // Replace ? in endpoint with & to append to existing query string
        url = this.baseUrl + endpoint.replace('?', '&');
      } else {
        url = this.baseUrl + endpoint;
      }
      const headers = {
        'Content-Type': 'application/json',
        ...(this.nonce ? { 'X-WP-Nonce': this.nonce } : {}),
      };

      const response = await fetch(url, {
        ...options,
        headers: { ...headers, ...options.headers },
      });

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      return response.json();
    },

    getProjects(params = {}) {
      const query = new URLSearchParams();
      if (params.tags?.length) {
        params.tags.forEach(tag => query.append('tags[]', tag));
      }
      if (params.hidden_categories?.length) {
        params.hidden_categories.forEach(cat => query.append('hidden_categories[]', cat));
      }
      if (params.page) query.set('page', params.page);
      if (params.per_page) query.set('per_page', params.per_page);

      return this.fetch(`search?${query.toString()}`);
    },
  };

  /**
   * Carousel Controller
   */
  class Carousel {
    constructor(element) {
      this.element = element;
      this.track = element.querySelector('.waldritter-carousel__slides');
      this.slides = Array.from(element.querySelectorAll('.waldritter-carousel__slide'));
      this.prevBtn = element.querySelector('.waldritter-carousel__button--prev');
      this.nextBtn = element.querySelector('.waldritter-carousel__button--next');
      this.dots = Array.from(element.querySelectorAll('.waldritter-carousel__dot'));

      this.currentIndex = 0;
      this.autoScroll = element.dataset.autoScroll === 'true';
      this.interval = parseInt(element.dataset.interval, 10) || 5000;
      this.autoScrollTimer = null;
      this.itemsPerView = this.getItemsPerView();

      this.init();
    }

    init() {
      if (this.slides.length <= 1) return;

      // Event listeners
      this.prevBtn?.addEventListener('click', () => this.prev());
      this.nextBtn?.addEventListener('click', () => this.next());
      this.dots.forEach((dot, index) => {
        dot.addEventListener('click', () => this.goTo(index));
      });

      // Keyboard navigation
      this.element.addEventListener('keydown', (e) => {
        if (e.key === 'ArrowLeft') this.prev();
        if (e.key === 'ArrowRight') this.next();
      });

      // Touch/swipe support
      let touchStartX = 0;
      let touchEndX = 0;

      this.track.addEventListener('touchstart', (e) => {
        touchStartX = e.changedTouches[0].screenX;
        this.stopAutoScroll();
      }, { passive: true });

      this.track.addEventListener('touchend', (e) => {
        touchEndX = e.changedTouches[0].screenX;
        const diff = touchStartX - touchEndX;
        if (Math.abs(diff) > 50) {
          if (diff > 0) this.next();
          else this.prev();
        }
        this.startAutoScroll();
      }, { passive: true });

      // Pause on hover
      this.element.addEventListener('mouseenter', () => this.stopAutoScroll());
      this.element.addEventListener('mouseleave', () => this.startAutoScroll());

      // Pause on focus
      this.element.addEventListener('focusin', () => this.stopAutoScroll());
      this.element.addEventListener('focusout', () => this.startAutoScroll());

      // Handle resize
      window.addEventListener('resize', () => {
        this.itemsPerView = this.getItemsPerView();
        this.updatePosition();
      });

      // Respect reduced motion preference
      if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
        this.autoScroll = false;
      }

      // Start auto scroll
      this.startAutoScroll();
      this.updateButtons();
    }

    getItemsPerView() {
      const width = window.innerWidth;
      if (width >= 1024) return 3;
      if (width >= 768) return 2;
      return 1;
    }

    getMaxIndex() {
      return Math.max(0, this.slides.length - this.itemsPerView);
    }

    prev() {
      this.goTo(Math.max(0, this.currentIndex - 1));
    }

    next() {
      const maxIndex = this.getMaxIndex();
      this.goTo(Math.min(maxIndex, this.currentIndex + 1));
    }

    goTo(index) {
      const maxIndex = this.getMaxIndex();
      this.currentIndex = Math.max(0, Math.min(maxIndex, index));
      this.updatePosition();
      this.updateButtons();
      this.updateDots();
    }

    updatePosition() {
      const slideWidth = 100 / this.itemsPerView;
      const offset = this.currentIndex * slideWidth;
      this.track.style.transform = `translateX(-${offset}%)`;
    }

    updateButtons() {
      const maxIndex = this.getMaxIndex();
      if (this.prevBtn) this.prevBtn.disabled = this.currentIndex === 0;
      if (this.nextBtn) this.nextBtn.disabled = this.currentIndex >= maxIndex;
    }

    updateDots() {
      this.dots.forEach((dot, index) => {
        const isActive = index === this.currentIndex;
        dot.classList.toggle('waldritter-carousel__dot--active', isActive);
        dot.setAttribute('aria-selected', isActive ? 'true' : 'false');
      });
    }

    startAutoScroll() {
      if (!this.autoScroll || this.autoScrollTimer) return;

      this.autoScrollTimer = setInterval(() => {
        const maxIndex = this.getMaxIndex();
        if (this.currentIndex >= maxIndex) {
          this.goTo(0);
        } else {
          this.next();
        }
      }, this.interval);
    }

    stopAutoScroll() {
      if (this.autoScrollTimer) {
        clearInterval(this.autoScrollTimer);
        this.autoScrollTimer = null;
      }
    }
  }

  /**
   * Search Controller
   */
  class Search {
    constructor(element) {
      this.element = element;
      this.id = element.id;
      this.perPage = parseInt(element.dataset.perPage, 10) || 5;
      this.hiddenCategories = JSON.parse(element.dataset.hiddenCategories || '[]');
      this.initialTags = JSON.parse(element.dataset.initialTags || '[]');
      this.basePatternTags = JSON.parse(element.dataset.basePatternTags || '[]');
      this.showOccurrences = element.dataset.showOccurrences === 'true';

      this.currentPage = 1;
      this.totalPages = 1;
      this.selectedTags = [...this.initialTags];

      // Elements
      this.filterButton = element.querySelector('.waldritter-search__filter-button');
      this.filterCount = element.querySelector('.waldritter-search__filter-count');
      this.drawer = element.querySelector('.waldritter-search__drawer');
      this.drawerClose = element.querySelector('.waldritter-search__drawer-close');
      this.backdrop = element.querySelector('.waldritter-search__drawer-backdrop');
      this.clearButton = element.querySelector('.waldritter-search__clear-button');
      this.applyButton = element.querySelector('.waldritter-search__apply-button');
      this.checkboxes = element.querySelectorAll('.waldritter-search__tag-checkbox');
      this.list = element.querySelector('.waldritter-search__list');
      this.resultsCount = element.querySelector('.waldritter-search__results-count');
      this.currentPageEl = element.querySelector('.waldritter-search__current-page');
      this.totalPagesEl = element.querySelector('.waldritter-search__total-pages');
      this.prevButton = element.querySelector('.waldritter-search__page-button--prev');
      this.nextButton = element.querySelector('.waldritter-search__page-button--next');

      this.init();
    }

    init() {
      // Drawer controls
      this.filterButton?.addEventListener('click', () => this.openDrawer());
      this.drawerClose?.addEventListener('click', () => this.closeDrawer());
      this.backdrop?.addEventListener('click', () => this.closeDrawer());

      // Keyboard close
      document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && this.isDrawerOpen()) {
          this.closeDrawer();
        }
      });

      // Filter buttons
      this.clearButton?.addEventListener('click', () => this.clearFilters());
      this.applyButton?.addEventListener('click', () => this.applyFilters());

      // Pagination
      this.prevButton?.addEventListener('click', () => this.goToPage(this.currentPage - 1));
      this.nextButton?.addEventListener('click', () => this.goToPage(this.currentPage + 1));

      // Occurrences toggles
      this.element.addEventListener('click', (e) => {
        const toggle = e.target.closest('.waldritter-project-card__occurrences-toggle');
        if (toggle) {
          this.toggleOccurrences(toggle);
        }
      });

      // Update filter count
      this.updateFilterCount();
    }

    isDrawerOpen() {
      return this.drawer?.getAttribute('aria-hidden') === 'false';
    }

    openDrawer() {
      this.drawer?.setAttribute('aria-hidden', 'false');
      this.filterButton?.setAttribute('aria-expanded', 'true');
      document.body.style.overflow = 'hidden';

      // Focus first checkbox
      const firstCheckbox = this.drawer?.querySelector('.waldritter-search__tag-checkbox');
      firstCheckbox?.focus();
    }

    closeDrawer() {
      this.drawer?.setAttribute('aria-hidden', 'true');
      this.filterButton?.setAttribute('aria-expanded', 'false');
      document.body.style.overflow = '';
      this.filterButton?.focus();
    }

    getSelectedTagsFromUI() {
      const tags = [];
      this.checkboxes.forEach((checkbox) => {
        if (checkbox.checked) {
          tags.push(checkbox.value);
        }
      });
      return tags;
    }

    clearFilters() {
      this.checkboxes.forEach((checkbox) => {
        checkbox.checked = false;
      });
    }

    async applyFilters() {
      this.selectedTags = this.getSelectedTagsFromUI();
      this.currentPage = 1;
      this.closeDrawer();
      this.updateFilterCount();
      await this.loadProjects();
    }

    updateFilterCount() {
      const count = this.selectedTags.length;
      if (this.filterCount) {
        this.filterCount.textContent = count > 0 ? count : '';
      }
    }

    async loadProjects() {
      this.element.classList.add('waldritter-search--loading');

      try {
        // Merge base pattern tags with user-selected tags (AND logic)
        const effectiveTags = [...this.basePatternTags, ...this.selectedTags];
        const data = await api.getProjects({
          tags: effectiveTags,
          hidden_categories: this.hiddenCategories,
          page: this.currentPage,
          per_page: this.perPage,
        });

        this.totalPages = data.pagination.total_pages;
        this.renderProjects(data.projects);
        this.updatePagination(data.pagination);
        this.updateResultsCount(data.pagination.total);
      } catch (error) {
        console.error('Failed to load projects:', error);
        this.list.innerHTML = '<div class="waldritter-search__empty"><p>Error loading projects. Please try again.</p></div>';
      } finally {
        this.element.classList.remove('waldritter-search--loading');
      }
    }

    renderProjects(projects) {
      if (!projects.length) {
        this.list.innerHTML = '<div class="waldritter-search__empty"><p>No projects found matching your criteria.</p></div>';
        return;
      }

      this.list.innerHTML = projects.map(project => this.renderProjectCard(project)).join('');
    }

    renderProjectCard(project) {
      const hasImage = !!project.imageUrl;
      const hasHomepage = !!project.homepage;

      // Get upcoming occurrences (up to 3)
      const upcomingOccurrences = dateHelper.getUpcoming(project.occurrences, 3);
      const totalUpcoming = project.occurrences
        ? project.occurrences.filter(o => dateHelper.isFuture(o.startDate)).length
        : 0;

      // Group tags by category
      const tagsByCategory = {};
      (project.tags || []).forEach(tag => {
        const category = tag.category?.title || 'Other';
        if (!tagsByCategory[category]) {
          tagsByCategory[category] = [];
        }
        tagsByCategory[category].push(tag.title);
      });

      const tagsHtml = Object.entries(tagsByCategory).map(([category, tags]) => `
        <div class="waldritter-project-card__tag-group">
          <span class="waldritter-project-card__tag-category">${this.escapeHtml(category)}:</span>
          <span class="waldritter-project-card__tag-list">${this.escapeHtml(tags.join(', '))}</span>
        </div>
      `).join('');

      const imageHtml = hasImage ? `
        <div class="waldritter-project-card__image">
          ${hasHomepage ? `<a href="${this.escapeHtml(project.homepage)}" target="_blank" rel="noopener noreferrer">` : ''}
            <img src="${this.escapeHtml(project.imageUrl)}" alt="${this.escapeHtml(project.title)}" loading="lazy" />
          ${hasHomepage ? '</a>' : ''}
        </div>
      ` : '';

      const linkHtml = hasHomepage ? `
        <a href="${this.escapeHtml(project.homepage)}" class="waldritter-project-card__link" target="_blank" rel="noopener noreferrer">
          Homepage besuchen
          <span class="waldritter-project-card__external-icon" aria-hidden="true">
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M10.5 1.5L1.5 10.5M10.5 1.5H4.5M10.5 1.5V7.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </span>
        </a>
      ` : '';

      // Build dates HTML
      let datesHtml = '';
      if (upcomingOccurrences.length > 0) {
        const datesListHtml = upcomingOccurrences.map((occ, index) => {
          const dateStr = dateHelper.formatRange(occ.startDate, occ.endDate);
          const relative = index === 0 ? dateHelper.getRelative(occ.startDate) : null;
          return `
            <div class="waldritter-project-card__date${index === 0 ? ' waldritter-project-card__date--next' : ''}">
              <span class="waldritter-project-card__date-value">${this.escapeHtml(dateStr)}</span>
              ${relative ? `<span class="waldritter-project-card__date-relative">${this.escapeHtml(relative)}</span>` : ''}
            </div>
          `;
        }).join('');

        const moreHtml = totalUpcoming > 3
          ? `<div class="waldritter-project-card__date waldritter-project-card__date--more">+ ${totalUpcoming - 3} weitere Termine</div>`
          : '';

        datesHtml = `
          <div class="waldritter-project-card__dates">
            <div class="waldritter-project-card__dates-icon" aria-hidden="true">
              <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="2" y="3" width="12" height="11" rx="1" stroke="currentColor" stroke-width="1.5"/>
                <path d="M2 6H14" stroke="currentColor" stroke-width="1.5"/>
                <path d="M5 1V4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                <path d="M11 1V4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
              </svg>
            </div>
            <div class="waldritter-project-card__dates-list">
              ${datesListHtml}
              ${moreHtml}
            </div>
          </div>
        `;
      }

      return `
        <article class="waldritter-project-card" role="listitem" data-project-id="${project.id}">
          <div class="waldritter-project-card__inner">
            ${imageHtml}
            <div class="waldritter-project-card__content">
              <h3 class="waldritter-project-card__title">${this.escapeHtml(project.title)}</h3>
              ${datesHtml}
              ${project.description ? `<div class="waldritter-project-card__description">${markdown.render(project.description)}</div>` : ''}
              ${tagsHtml ? `<div class="waldritter-project-card__tags">${tagsHtml}</div>` : ''}
              <div class="waldritter-project-card__actions">${linkHtml}</div>
            </div>
          </div>
        </article>
      `;
    }

    escapeHtml(text) {
      const div = document.createElement('div');
      div.textContent = text;
      return div.innerHTML;
    }

    updatePagination(pagination) {
      if (this.currentPageEl) {
        this.currentPageEl.textContent = pagination.page;
      }
      if (this.totalPagesEl) {
        this.totalPagesEl.textContent = pagination.total_pages;
      }
      if (this.prevButton) {
        this.prevButton.disabled = pagination.page <= 1;
      }
      if (this.nextButton) {
        this.nextButton.disabled = pagination.page >= pagination.total_pages;
      }
    }

    updateResultsCount(total) {
      if (this.resultsCount) {
        const text = total === 1 ? '1 Project' : `${total} Projects`;
        this.resultsCount.textContent = text;
      }
    }

    async goToPage(page) {
      if (page < 1 || page > this.totalPages) return;
      this.currentPage = page;
      await this.loadProjects();

      // Scroll to top of list
      this.element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }

    toggleOccurrences(toggle) {
      const expanded = toggle.getAttribute('aria-expanded') === 'true';
      const targetId = toggle.getAttribute('aria-controls');
      const target = document.getElementById(targetId);

      toggle.setAttribute('aria-expanded', !expanded);
      if (target) {
        target.hidden = expanded;
      }
    }
  }

  /**
   * Initialize all widgets on page load
   */
  function init() {
    // Initialize markdown renderer
    markdown.init();

    // Initialize carousels
    document.querySelectorAll('.waldritter-carousel').forEach((el) => {
      new Carousel(el);
    });

    // Initialize search widgets
    document.querySelectorAll('.waldritter-search').forEach((el) => {
      new Search(el);
    });
  }

  // Run on DOM ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
