/**
 * Waldritter Project DB - Gutenberg Blocks
 */

(function(wp) {
  'use strict';

  const { registerBlockType } = wp.blocks;
  const { InspectorControls } = wp.blockEditor;
  const {
    PanelBody,
    TextControl,
    ToggleControl,
    RangeControl,
    Spinner,
  } = wp.components;
  const { Fragment, useState, useEffect } = wp.element;
  const { __ } = wp.i18n;
  const apiFetch = wp.apiFetch;

  const BLOCK_ICON = wp.element.createElement('svg', {
    width: 24,
    height: 24,
    viewBox: '0 0 24 24',
    fill: 'none',
    xmlns: 'http://www.w3.org/2000/svg',
  }, wp.element.createElement('path', {
    d: 'M3 7C3 5.89543 3.89543 5 5 5H19C20.1046 5 21 5.89543 21 7V17C21 18.1046 20.1046 19 19 19H5C3.89543 19 3 18.1046 3 17V7Z',
    stroke: 'currentColor',
    strokeWidth: 2,
    strokeLinecap: 'round',
    strokeLinejoin: 'round',
  }), wp.element.createElement('path', {
    d: 'M8 12H16M8 16H12',
    stroke: 'currentColor',
    strokeWidth: 2,
    strokeLinecap: 'round',
    strokeLinejoin: 'round',
  }));

  /**
   * Tag Selector Component
   */
  function TagSelector({ selectedTags, onChange }) {
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
      apiFetch({ path: '/waldritter/v1/categories' })
        .then((data) => {
          setCategories(data);
          setLoading(false);
        })
        .catch((err) => {
          setError(err.message);
          setLoading(false);
        });
    }, []);

    if (loading) {
      return wp.element.createElement(Spinner);
    }

    if (error) {
      return wp.element.createElement('p', { className: 'waldritter-block-error' }, error);
    }

    const toggleTag = (tagTitle) => {
      const newTags = selectedTags.includes(tagTitle)
        ? selectedTags.filter(t => t !== tagTitle)
        : [...selectedTags, tagTitle];
      onChange(newTags);
    };

    return wp.element.createElement('div', { className: 'waldritter-tag-selector' },
      categories.map((category) =>
        wp.element.createElement('div', {
          key: category.id,
          className: 'waldritter-tag-selector__category',
        },
          wp.element.createElement('div', {
            className: 'waldritter-tag-selector__category-name',
          }, category.title),
          wp.element.createElement('div', {
            className: 'waldritter-tag-selector__tags',
          },
            (category.tags || []).map((tag) =>
              wp.element.createElement('button', {
                key: tag.id,
                type: 'button',
                className: `waldritter-tag-selector__tag ${selectedTags.includes(tag.title) ? 'waldritter-tag-selector__tag--selected' : ''}`,
                onClick: () => toggleTag(tag.title),
              }, tag.title)
            )
          )
        )
      )
    );
  }

  /**
   * Carousel Block
   */
  registerBlockType('waldritter/carousel', {
    title: __('Waldritter Carousel', 'waldritter-project-db'),
    description: __('Carousel displaying featured projects with tag filtering.', 'waldritter-project-db'),
    category: 'widgets',
    icon: BLOCK_ICON,
    keywords: [__('carousel', 'waldritter-project-db'), __('slider', 'waldritter-project-db'), __('projects', 'waldritter-project-db')],
    supports: {
      html: false,
      align: ['wide', 'full'],
    },

    edit: function(props) {
      const { attributes, setAttributes } = props;
      const { tags, projectIds, hiddenCategories, autoScroll, interval } = attributes;

      return wp.element.createElement(Fragment, null,
        wp.element.createElement(InspectorControls, null,
          wp.element.createElement(PanelBody, {
            title: __('Filter by Tags', 'waldritter-project-db'),
            initialOpen: true,
          },
            wp.element.createElement('p', { style: { color: '#757575', fontSize: '12px' } },
              __('Select tags to filter projects (AND logic)', 'waldritter-project-db')
            ),
            wp.element.createElement(TagSelector, {
              selectedTags: tags,
              onChange: (newTags) => setAttributes({ tags: newTags }),
            })
          ),
          wp.element.createElement(PanelBody, {
            title: __('Specific Projects', 'waldritter-project-db'),
            initialOpen: false,
          },
            wp.element.createElement(TextControl, {
              label: __('Project IDs', 'waldritter-project-db'),
              help: __('Comma-separated list of project IDs', 'waldritter-project-db'),
              value: (projectIds || []).join(', '),
              onChange: (value) => setAttributes({
                projectIds: value.split(',').map(id => id.trim()).filter(Boolean),
              }),
            })
          ),
          wp.element.createElement(PanelBody, {
            title: __('Hidden Categories', 'waldritter-project-db'),
            initialOpen: false,
          },
            wp.element.createElement(TextControl, {
              label: __('Category Names', 'waldritter-project-db'),
              help: __('Comma-separated list of category names to hide', 'waldritter-project-db'),
              value: (hiddenCategories || []).join(', '),
              onChange: (value) => setAttributes({
                hiddenCategories: value.split(',').map(c => c.trim()).filter(Boolean),
              }),
            })
          ),
          wp.element.createElement(PanelBody, {
            title: __('Animation', 'waldritter-project-db'),
            initialOpen: false,
          },
            wp.element.createElement(ToggleControl, {
              label: __('Auto Scroll', 'waldritter-project-db'),
              checked: autoScroll,
              onChange: (value) => setAttributes({ autoScroll: value }),
            }),
            autoScroll && wp.element.createElement(RangeControl, {
              label: __('Interval (ms)', 'waldritter-project-db'),
              value: interval,
              onChange: (value) => setAttributes({ interval: value }),
              min: 2000,
              max: 10000,
              step: 500,
            })
          )
        ),
        wp.element.createElement('div', { className: 'waldritter-carousel-preview' },
          [1, 2, 3].map((i) =>
            wp.element.createElement('div', { key: i, className: 'waldritter-carousel-preview__card' },
              wp.element.createElement('div', { className: 'waldritter-carousel-preview__image' }),
              wp.element.createElement('div', { className: 'waldritter-carousel-preview__content' },
                wp.element.createElement('div', { className: 'waldritter-carousel-preview__title' },
                  __('Project', 'waldritter-project-db') + ' ' + i
                )
              )
            )
          )
        )
      );
    },

    save: function() {
      return null; // Server-side rendered
    },
  });

  /**
   * Search Block
   */
  registerBlockType('waldritter/search', {
    title: __('Waldritter Search', 'waldritter-project-db'),
    description: __('Interactive project browser with tag filtering and pagination.', 'waldritter-project-db'),
    category: 'widgets',
    icon: BLOCK_ICON,
    keywords: [__('search', 'waldritter-project-db'), __('filter', 'waldritter-project-db'), __('projects', 'waldritter-project-db')],
    supports: {
      html: false,
      align: ['wide', 'full'],
    },

    edit: function(props) {
      const { attributes, setAttributes } = props;
      const { perPage, hiddenCategories, initialTags, basePatternTags, showOccurrences } = attributes;

      return wp.element.createElement(Fragment, null,
        wp.element.createElement(InspectorControls, null,
          wp.element.createElement(PanelBody, {
            title: __('Display Settings', 'waldritter-project-db'),
            initialOpen: true,
          },
            wp.element.createElement(RangeControl, {
              label: __('Projects per Page', 'waldritter-project-db'),
              value: perPage,
              onChange: (value) => setAttributes({ perPage: value }),
              min: 1,
              max: 20,
            }),
            wp.element.createElement(ToggleControl, {
              label: __('Show Occurrences', 'waldritter-project-db'),
              help: __('Display expandable event dates', 'waldritter-project-db'),
              checked: showOccurrences,
              onChange: (value) => setAttributes({ showOccurrences: value }),
            })
          ),
          wp.element.createElement(PanelBody, {
            title: __('Initial Tags', 'waldritter-project-db'),
            initialOpen: false,
          },
            wp.element.createElement('p', { style: { color: '#757575', fontSize: '12px' } },
              __('Pre-select these tags on page load', 'waldritter-project-db')
            ),
            wp.element.createElement(TagSelector, {
              selectedTags: initialTags,
              onChange: (newTags) => setAttributes({ initialTags: newTags }),
            })
          ),
          wp.element.createElement(PanelBody, {
            title: __('Base Pattern (Always Applied)', 'waldritter-project-db'),
            initialOpen: false,
          },
            wp.element.createElement('p', { style: { color: '#757575', fontSize: '12px', marginBottom: '10px' } },
              __('These tags are always applied to searches and cannot be removed by users.', 'waldritter-project-db')
            ),
            wp.element.createElement(TagSelector, {
              selectedTags: basePatternTags || [],
              onChange: (newTags) => setAttributes({ basePatternTags: newTags }),
            })
          ),
          wp.element.createElement(PanelBody, {
            title: __('Hidden Categories', 'waldritter-project-db'),
            initialOpen: false,
          },
            wp.element.createElement(TextControl, {
              label: __('Category Names', 'waldritter-project-db'),
              help: __('Comma-separated list of category names to hide from filter', 'waldritter-project-db'),
              value: (hiddenCategories || []).join(', '),
              onChange: (value) => setAttributes({
                hiddenCategories: value.split(',').map(c => c.trim()).filter(Boolean),
              }),
            })
          )
        ),
        wp.element.createElement('div', { className: 'waldritter-search-preview' },
          wp.element.createElement('div', { className: 'waldritter-search-preview__filter' },
            wp.element.createElement('span', null, __('Filter', 'waldritter-project-db')),
            initialTags.length > 0 && wp.element.createElement('span', null, `(${initialTags.length})`)
          ),
          wp.element.createElement('div', { className: 'waldritter-search-preview__list' },
            [1, 2].map((i) =>
              wp.element.createElement('div', { key: i, className: 'waldritter-search-preview__card' },
                wp.element.createElement('div', { className: 'waldritter-search-preview__card-image' }),
                wp.element.createElement('div', { className: 'waldritter-search-preview__card-content' },
                  wp.element.createElement('div', { className: 'waldritter-search-preview__card-title' },
                    __('Project', 'waldritter-project-db') + ' ' + i
                  ),
                  wp.element.createElement('p', { className: 'waldritter-search-preview__card-description' },
                    __('Project description...', 'waldritter-project-db')
                  )
                )
              )
            )
          )
        )
      );
    },

    save: function() {
      return null; // Server-side rendered
    },
  });

})(window.wp);
