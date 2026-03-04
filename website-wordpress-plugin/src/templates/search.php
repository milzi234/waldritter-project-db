<?php
/**
 * Search widget template
 *
 * @package Waldritter\ProjectDB
 *
 * @var array $attributes Widget attributes
 * @var array $categories Categories for filter
 * @var array $projects Initial projects
 * @var array $hidden_categories Categories to hide
 * @var array $initial_tags Initially selected tags
 * @var array $base_pattern_tags Base pattern tags (always applied, hidden from UI)
 */

declare(strict_types=1);

defined('ABSPATH') || exit;

$widget_id = esc_attr($attributes['id']);
$per_page = (int) $attributes['per_page'];
$show_occurrences = $attributes['show_occurrences'];

// Calculate initial pagination
$total = count($projects);
$total_pages = (int) ceil($total / $per_page);
$initial_projects = array_slice($projects, 0, $per_page);
?>
<div
    id="<?php echo $widget_id; ?>"
    class="waldritter-search"
    data-per-page="<?php echo esc_attr($per_page); ?>"
    data-hidden-categories="<?php echo esc_attr(wp_json_encode($hidden_categories)); ?>"
    data-initial-tags="<?php echo esc_attr(wp_json_encode($initial_tags)); ?>"
    data-base-pattern-tags="<?php echo esc_attr(wp_json_encode($base_pattern_tags ?? [])); ?>"
    data-show-occurrences="<?php echo $show_occurrences ? 'true' : 'false'; ?>"
>
    <!-- Filter Button & Drawer -->
    <div class="waldritter-search__filter-section">
        <button
            type="button"
            class="waldritter-search__filter-button"
            aria-expanded="false"
            aria-controls="<?php echo $widget_id; ?>-drawer"
        >
            <span class="waldritter-search__filter-icon" aria-hidden="true">
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 10H15M2.5 5H17.5M7.5 15H12.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
            </span>
            <span class="waldritter-search__filter-text">
                <?php esc_html_e('Filter', 'waldritter-project-db'); ?>
            </span>
            <span class="waldritter-search__filter-count" aria-live="polite"></span>
        </button>

        <!-- Filter Drawer -->
        <div
            id="<?php echo $widget_id; ?>-drawer"
            class="waldritter-search__drawer"
            aria-hidden="true"
        >
            <div class="waldritter-search__drawer-header">
                <h3 class="waldritter-search__drawer-title">
                    <?php esc_html_e('Filter by Tags', 'waldritter-project-db'); ?>
                </h3>
                <button
                    type="button"
                    class="waldritter-search__drawer-close"
                    aria-label="<?php esc_attr_e('Close filter', 'waldritter-project-db'); ?>"
                >
                    <span aria-hidden="true">&times;</span>
                </button>
            </div>

            <div class="waldritter-search__drawer-content">
                <?php foreach ($categories as $category): ?>
                    <fieldset class="waldritter-search__category">
                        <legend class="waldritter-search__category-title">
                            <?php echo esc_html($category['title']); ?>
                        </legend>
                        <div class="waldritter-search__tags">
                            <?php foreach ($category['tags'] ?? [] as $tag): ?>
                                <?php
                                $is_checked = in_array($tag['title'], $initial_tags, true);
                                $tag_id = $widget_id . '-tag-' . esc_attr($tag['id']);
                                ?>
                                <label class="waldritter-search__tag-label" for="<?php echo $tag_id; ?>">
                                    <input
                                        type="checkbox"
                                        id="<?php echo $tag_id; ?>"
                                        class="waldritter-search__tag-checkbox"
                                        name="tags[]"
                                        value="<?php echo esc_attr($tag['title']); ?>"
                                        data-tag-id="<?php echo esc_attr($tag['id']); ?>"
                                        <?php checked($is_checked); ?>
                                    />
                                    <span class="waldritter-search__tag-text">
                                        <?php echo esc_html($tag['title']); ?>
                                    </span>
                                </label>
                            <?php endforeach; ?>
                        </div>
                    </fieldset>
                <?php endforeach; ?>
            </div>

            <div class="waldritter-search__drawer-footer">
                <button
                    type="button"
                    class="waldritter-search__clear-button"
                >
                    <?php esc_html_e('Clear All', 'waldritter-project-db'); ?>
                </button>
                <button
                    type="button"
                    class="waldritter-search__apply-button"
                >
                    <?php esc_html_e('Apply Filters', 'waldritter-project-db'); ?>
                </button>
            </div>
        </div>
        <div class="waldritter-search__drawer-backdrop" aria-hidden="true"></div>
    </div>

    <!-- Results Count -->
    <div class="waldritter-search__results-info" aria-live="polite">
        <span class="waldritter-search__results-count">
            <?php echo esc_html(sprintf(
                /* translators: %d: number of projects */
                _n('%d Project', '%d Projects', $total, 'waldritter-project-db'),
                $total
            )); ?>
        </span>
    </div>

    <!-- Project List -->
    <div class="waldritter-search__list" role="list">
        <?php if (empty($initial_projects)): ?>
            <div class="waldritter-search__empty">
                <p><?php esc_html_e('No projects found matching your criteria.', 'waldritter-project-db'); ?></p>
            </div>
        <?php else: ?>
            <?php foreach ($initial_projects as $project): ?>
                <?php include __DIR__ . '/partials/project-card.php'; ?>
            <?php endforeach; ?>
        <?php endif; ?>
    </div>

    <!-- Loading State -->
    <div class="waldritter-search__loading" aria-hidden="true">
        <div class="waldritter-search__spinner"></div>
        <span class="screen-reader-text"><?php esc_html_e('Loading projects...', 'waldritter-project-db'); ?></span>
    </div>

    <!-- Pagination -->
    <?php if ($total_pages > 1): ?>
        <nav
            class="waldritter-search__pagination"
            aria-label="<?php esc_attr_e('Project pagination', 'waldritter-project-db'); ?>"
        >
            <button
                type="button"
                class="waldritter-search__page-button waldritter-search__page-button--prev"
                disabled
                aria-label="<?php esc_attr_e('Previous page', 'waldritter-project-db'); ?>"
            >
                <span aria-hidden="true">&laquo;</span>
                <span class="waldritter-search__page-button-text"><?php esc_html_e('Previous', 'waldritter-project-db'); ?></span>
            </button>

            <div class="waldritter-search__page-info">
                <span class="waldritter-search__current-page">1</span>
                <span class="waldritter-search__page-separator">/</span>
                <span class="waldritter-search__total-pages"><?php echo esc_html($total_pages); ?></span>
            </div>

            <button
                type="button"
                class="waldritter-search__page-button waldritter-search__page-button--next"
                <?php disabled($total_pages <= 1); ?>
                aria-label="<?php esc_attr_e('Next page', 'waldritter-project-db'); ?>"
            >
                <span class="waldritter-search__page-button-text"><?php esc_html_e('Next', 'waldritter-project-db'); ?></span>
                <span aria-hidden="true">&raquo;</span>
            </button>
        </nav>
    <?php endif; ?>
</div>
