import { distanceInWordsStrict } from 'date-fns'
import { isAfter } from 'date-fns'
import { numberWithCommas, pluralize } from '../../../packages/webapp/src/util/strings'

/**
 * Returns "N users" (properly pluralized and with commas added to N as needed).
 */
export function formatUserCount(userCount: number, hyphenate?: boolean): string {
    if (hyphenate) {
        return `${numberWithCommas(userCount)}-user`
    }
    return `${numberWithCommas(userCount)} ${pluralize('user', userCount)}`
}

/**
 * Reports whether {@link expiresAt} is in the past.
 */
export function isProductLicenseExpired(expiresAt: string | number): boolean {
    return !isAfter(expiresAt, Date.now())
}

/**
 * Returns "T remaining" or "T ago" for an expiration date.
 */
export function formatRelativeExpirationDate(date: string | number): string {
    return `${distanceInWordsStrict(date, Date.now())} ${isProductLicenseExpired(date) ? 'ago' : 'remaining'}`
}

/**
 * Returns a mailto:sales@sourcegraph.com URL with an optional subject.
 */
export function mailtoSales(args: { subject?: string }): string {
    return `mailto:sales@sourcegraph.com?subject=${encodeURIComponent(args.subject || '')}`
}
