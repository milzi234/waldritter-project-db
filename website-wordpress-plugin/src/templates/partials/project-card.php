<?php
/**
 * Project card partial template
 *
 * @package Waldritter\ProjectDB
 *
 * @var array $project Project data
 * @var bool $show_occurrences Whether to show occurrences
 */

declare(strict_types=1);

defined('ABSPATH') || exit;

use Waldritter\ProjectDB\DateHelper;

$project_id = esc_attr($project['id']);
$has_image = !empty($project['imageUrl']);
$has_homepage = !empty($project['homepage']);

// Group tags by category
$tags_by_category = [];
foreach ($project['tags'] ?? [] as $tag) {
    $category_title = $tag['category']['title'] ?? __('Other', 'waldritter-project-db');
    if (!isset($tags_by_category[$category_title])) {
        $tags_by_category[$category_title] = [];
    }
    $tags_by_category[$category_title][] = $tag['title'];
}

// Get upcoming occurrences (filter to future dates only)
$upcoming_occurrences = [];
$now = new DateTime();
foreach ($project['occurrences'] ?? [] as $occurrence) {
    if (!empty($occurrence['startDate'])) {
        try {
            $start = new DateTime($occurrence['startDate']);
            if ($start >= $now) {
                $upcoming_occurrences[] = $occurrence;
            }
        } catch (Exception $e) {
            // Skip invalid dates
        }
    }
}
// Limit to first 3 upcoming dates for display
$display_occurrences = array_slice($upcoming_occurrences, 0, 3);
$has_upcoming_dates = !empty($display_occurrences);
?>
<article class="waldritter-project-card" role="listitem" data-project-id="<?php echo $project_id; ?>">
    <div class="waldritter-project-card__inner">
        <?php if ($has_image): ?>
            <div class="waldritter-project-card__image">
                <?php if ($has_homepage): ?>
                    <a href="<?php echo esc_url($project['homepage']); ?>" target="_blank" rel="noopener noreferrer">
                        <img
                            src="<?php echo esc_url($project['imageUrl']); ?>"
                            alt="<?php echo esc_attr($project['title']); ?>"
                            loading="lazy"
                        />
                    </a>
                <?php else: ?>
                    <img
                        src="<?php echo esc_url($project['imageUrl']); ?>"
                        alt="<?php echo esc_attr($project['title']); ?>"
                        loading="lazy"
                    />
                <?php endif; ?>
            </div>
        <?php endif; ?>

        <div class="waldritter-project-card__content">
            <h3 class="waldritter-project-card__title">
                <?php echo esc_html($project['title']); ?>
            </h3>

            <?php if ($has_upcoming_dates): ?>
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
                        <?php foreach ($display_occurrences as $index => $occurrence): ?>
                            <?php
                            $date_str = DateHelper::formatRange(
                                $occurrence['startDate'],
                                $occurrence['endDate'] ?? null
                            );
                            $relative = DateHelper::getRelative($occurrence['startDate']);
                            ?>
                            <div class="waldritter-project-card__date<?php echo $index === 0 ? ' waldritter-project-card__date--next' : ''; ?>">
                                <span class="waldritter-project-card__date-value"><?php echo esc_html($date_str); ?></span>
                                <?php if ($index === 0 && $relative): ?>
                                    <span class="waldritter-project-card__date-relative"><?php echo esc_html($relative); ?></span>
                                <?php endif; ?>
                            </div>
                        <?php endforeach; ?>
                        <?php if (count($upcoming_occurrences) > 3): ?>
                            <div class="waldritter-project-card__date waldritter-project-card__date--more">
                                <?php
                                $more_count = count($upcoming_occurrences) - 3;
                                echo esc_html(sprintf(
                                    $more_count === 1 ? '+ %d weiterer Termin' : '+ %d weitere Termine',
                                    $more_count
                                )); ?>
                            </div>
                        <?php endif; ?>
                    </div>
                </div>
            <?php endif; ?>

            <?php if (!empty($project['description'])): ?>
                <div class="waldritter-project-card__description">
                    <?php echo \Waldritter\ProjectDB\MarkdownRenderer::render($project['description']); ?>
                </div>
            <?php endif; ?>

            <?php if (!empty($tags_by_category)): ?>
                <div class="waldritter-project-card__tags">
                    <?php foreach ($tags_by_category as $category => $tags): ?>
                        <div class="waldritter-project-card__tag-group">
                            <span class="waldritter-project-card__tag-category">
                                <?php echo esc_html($category); ?>:
                            </span>
                            <span class="waldritter-project-card__tag-list">
                                <?php echo esc_html(implode(', ', $tags)); ?>
                            </span>
                        </div>
                    <?php endforeach; ?>
                </div>
            <?php endif; ?>

            <div class="waldritter-project-card__actions">
                <?php if ($has_homepage): ?>
                    <a
                        href="<?php echo esc_url($project['homepage']); ?>"
                        class="waldritter-project-card__link"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        <?php esc_html_e('Homepage besuchen', 'waldritter-project-db'); ?>
                        <span class="waldritter-project-card__external-icon" aria-hidden="true">
                            <svg width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <path d="M10.5 1.5L1.5 10.5M10.5 1.5H4.5M10.5 1.5V7.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                        </span>
                    </a>
                <?php endif; ?>
            </div>

        </div>
    </div>
</article>
