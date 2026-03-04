<?php
/**
 * Carousel widget template
 *
 * @package Waldritter\ProjectDB
 *
 * @var array $attributes Widget attributes
 * @var array $projects Projects from API
 * @var array $highlights Manual highlight items
 */

declare(strict_types=1);

defined('ABSPATH') || exit;

use Waldritter\ProjectDB\DateHelper;

// Combine manual highlights with projects
$items = [];

// Add manual highlights first
foreach ($highlights as $highlight) {
    $items[] = [
        'type' => 'highlight',
        'title' => $highlight['title'] ?? '',
        'description' => $highlight['description'] ?? '',
        'image' => $highlight['image']['url'] ?? '',
        'link' => $highlight['link'] ?? '',
        'link_text' => $highlight['linkText'] ?? __('Homepage besuchen', 'waldritter-project-db'),
    ];
}

// Add projects
foreach ($projects as $project) {
    $items[] = [
        'type' => 'project',
        'id' => $project['id'],
        'title' => $project['title'] ?? '',
        'description' => $project['description'] ?? '',
        'image' => $project['imageUrl'] ?? '',
        'link' => $project['homepage'] ?? '',
        'link_text' => __('Homepage besuchen', 'waldritter-project-db'),
        'tags' => $project['tags'] ?? [],
        'next_occurrence' => $project['nextOccurrence'] ?? null,
    ];
}

// If no items, show placeholder
if (empty($items)) {
    $items = [
        [
            'type' => 'placeholder',
            'title' => __('Keine Projekte gefunden', 'waldritter-project-db'),
            'description' => __('Es wurden keine Projekte mit den angegebenen Filtern gefunden.', 'waldritter-project-db'),
            'image' => '',
            'link' => '',
        ],
    ];
}

$carousel_id = esc_attr($attributes['id']);
$auto_scroll = $attributes['auto_scroll'] ? 'true' : 'false';
$interval = (int) $attributes['interval'];
?>
<div
    id="<?php echo $carousel_id; ?>"
    class="waldritter-carousel"
    data-auto-scroll="<?php echo esc_attr($auto_scroll); ?>"
    data-interval="<?php echo esc_attr($interval); ?>"
    role="region"
    aria-label="<?php esc_attr_e('Project Carousel', 'waldritter-project-db'); ?>"
>
    <div class="waldritter-carousel__track">
        <div class="waldritter-carousel__slides">
            <?php foreach ($items as $index => $item): ?>
                <div
                    class="waldritter-carousel__slide"
                    role="group"
                    aria-roledescription="slide"
                    aria-label="<?php echo esc_attr(sprintf(
                        /* translators: 1: current slide number, 2: total slides */
                        __('%1$d of %2$d', 'waldritter-project-db'),
                        $index + 1,
                        count($items)
                    )); ?>"
                >
                    <article class="waldritter-carousel__card">
                        <?php if (!empty($item['image'])): ?>
                            <div class="waldritter-carousel__image">
                                <?php if (!empty($item['link'])): ?>
                                    <a href="<?php echo esc_url($item['link']); ?>" tabindex="-1">
                                        <img
                                            src="<?php echo esc_url($item['image']); ?>"
                                            alt="<?php echo esc_attr($item['title']); ?>"
                                            loading="lazy"
                                        />
                                    </a>
                                <?php else: ?>
                                    <img
                                        src="<?php echo esc_url($item['image']); ?>"
                                        alt="<?php echo esc_attr($item['title']); ?>"
                                        loading="lazy"
                                    />
                                <?php endif; ?>
                            </div>
                        <?php endif; ?>

                        <div class="waldritter-carousel__content">
                            <h3 class="waldritter-carousel__title">
                                <?php echo esc_html($item['title']); ?>
                            </h3>

                            <?php if (!empty($item['next_occurrence'])): ?>
                                <?php
                                $next_date = DateHelper::formatRange(
                                    $item['next_occurrence']['startDateTime'],
                                    $item['next_occurrence']['endDateTime'] ?? null
                                );
                                $relative = DateHelper::getRelative($item['next_occurrence']['startDateTime']);
                                ?>
                                <div class="waldritter-carousel__date">
                                    <div class="waldritter-carousel__date-row">
                                        <span class="waldritter-carousel__date-icon" aria-hidden="true">
                                            <svg width="14" height="14" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                                                <rect x="2" y="3" width="12" height="11" rx="1" stroke="currentColor" stroke-width="1.5"/>
                                                <path d="M2 6H14" stroke="currentColor" stroke-width="1.5"/>
                                                <path d="M5 1V4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                                                <path d="M11 1V4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                                            </svg>
                                        </span>
                                        <span class="waldritter-carousel__date-value"><?php echo esc_html($next_date); ?></span>
                                    </div>
                                    <?php if ($relative): ?>
                                        <span class="waldritter-carousel__date-relative"><?php echo esc_html($relative); ?></span>
                                    <?php endif; ?>
                                </div>
                            <?php endif; ?>

                            <?php if (!empty($item['description'])): ?>
                                <p class="waldritter-carousel__description">
                                    <?php echo esc_html(\Waldritter\ProjectDB\MarkdownRenderer::renderExcerpt($item['description'], 20)); ?>
                                </p>
                            <?php endif; ?>

                            <?php if (!empty($item['tags'])): ?>
                                <div class="waldritter-carousel__tags">
                                    <?php foreach (array_slice($item['tags'], 0, 3) as $tag): ?>
                                        <span class="waldritter-carousel__tag">
                                            <?php echo esc_html($tag['title']); ?>
                                        </span>
                                    <?php endforeach; ?>
                                </div>
                            <?php endif; ?>

                            <?php if (!empty($item['link'])): ?>
                                <a
                                    href="<?php echo esc_url($item['link']); ?>"
                                    class="waldritter-carousel__link"
                                >
                                    <?php echo esc_html($item['link_text']); ?>
                                    <span class="waldritter-carousel__link-arrow" aria-hidden="true">&rarr;</span>
                                </a>
                            <?php endif; ?>
                        </div>
                    </article>
                </div>
            <?php endforeach; ?>
        </div>
    </div>

    <?php if (count($items) > 1): ?>
        <div class="waldritter-carousel__controls">
            <button
                type="button"
                class="waldritter-carousel__button waldritter-carousel__button--prev"
                aria-label="<?php esc_attr_e('Previous slide', 'waldritter-project-db'); ?>"
            >
                <span aria-hidden="true">&lsaquo;</span>
            </button>

            <div class="waldritter-carousel__dots" role="tablist">
                <?php foreach ($items as $index => $item): ?>
                    <button
                        type="button"
                        class="waldritter-carousel__dot <?php echo $index === 0 ? 'waldritter-carousel__dot--active' : ''; ?>"
                        role="tab"
                        aria-selected="<?php echo $index === 0 ? 'true' : 'false'; ?>"
                        aria-label="<?php echo esc_attr(sprintf(
                            /* translators: %d: slide number */
                            __('Go to slide %d', 'waldritter-project-db'),
                            $index + 1
                        )); ?>"
                        data-slide="<?php echo esc_attr($index); ?>"
                    ></button>
                <?php endforeach; ?>
            </div>

            <button
                type="button"
                class="waldritter-carousel__button waldritter-carousel__button--next"
                aria-label="<?php esc_attr_e('Next slide', 'waldritter-project-db'); ?>"
            >
                <span aria-hidden="true">&rsaquo;</span>
            </button>
        </div>
    <?php endif; ?>
</div>
