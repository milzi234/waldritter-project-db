<?php
/**
 * Date Helper
 *
 * @package Waldritter\ProjectDB
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB;

use DateTime;
use DateTimeZone;

/**
 * Helper class for formatting dates
 */
class DateHelper {
    /**
     * Format a date/time for display
     *
     * @param string|null $datetime ISO8601 datetime string
     * @param string $format PHP date format
     * @return string|null Formatted date or null
     */
    public static function format(?string $datetime, string $format = 'd.m.Y'): ?string {
        if (empty($datetime)) {
            return null;
        }

        try {
            $date = new DateTime($datetime);
            return $date->format($format);
        } catch (\Exception $e) {
            return null;
        }
    }

    /**
     * Format a date range for display
     *
     * @param string|null $start_datetime Start datetime
     * @param string|null $end_datetime End datetime
     * @return string|null Formatted range
     */
    public static function formatRange(?string $start_datetime, ?string $end_datetime = null): ?string {
        if (empty($start_datetime)) {
            return null;
        }

        try {
            $start = new DateTime($start_datetime);
            $start_date = $start->format('d.m.Y');
            $start_time = $start->format('H:i');

            if (empty($end_datetime)) {
                // Single datetime
                if ($start_time !== '00:00') {
                    return $start_date . ' ' . $start_time;
                }
                return $start_date;
            }

            $end = new DateTime($end_datetime);
            $end_date = $end->format('d.m.Y');
            $end_time = $end->format('H:i');

            // Same day
            if ($start_date === $end_date) {
                if ($start_time !== '00:00' || $end_time !== '00:00') {
                    return $start_date . ' ' . $start_time . ' - ' . $end_time;
                }
                return $start_date;
            }

            // Different days
            if ($start_time !== '00:00' || $end_time !== '00:00') {
                return $start_date . ' ' . $start_time . ' - ' . $end_date . ' ' . $end_time;
            }
            return $start_date . ' - ' . $end_date;
        } catch (\Exception $e) {
            return null;
        }
    }

    /**
     * Check if a datetime is in the future
     *
     * @param string|null $datetime ISO8601 datetime string
     * @return bool
     */
    public static function isFuture(?string $datetime): bool {
        if (empty($datetime)) {
            return false;
        }

        try {
            $date = new DateTime($datetime);
            $now = new DateTime();
            return $date > $now;
        } catch (\Exception $e) {
            return false;
        }
    }

    /**
     * Get a human-readable relative date
     *
     * @param string|null $datetime ISO8601 datetime string
     * @return string|null
     */
    public static function getRelative(?string $datetime): ?string {
        if (empty($datetime)) {
            return null;
        }

        try {
            $date = new DateTime($datetime);
            $now = new DateTime();
            $diff = $now->diff($date);

            if ($date < $now) {
                return null; // Past date
            }

            if ($diff->days === 0) {
                return 'Heute';
            }
            if ($diff->days === 1) {
                return 'Morgen';
            }
            if ($diff->days < 7) {
                return sprintf(
                    $diff->days === 1 ? 'In %d Tag' : 'In %d Tagen',
                    $diff->days
                );
            }
            if ($diff->days < 30) {
                $weeks = (int) floor($diff->days / 7);
                return sprintf(
                    $weeks === 1 ? 'In %d Woche' : 'In %d Wochen',
                    $weeks
                );
            }

            return null; // Too far in future for relative date
        } catch (\Exception $e) {
            return null;
        }
    }
}
