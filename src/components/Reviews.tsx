import { useEffect, useState } from 'react'
import { Star } from 'lucide-react'
import type { Review } from '../types'
import { storageService } from '../services/storage'
import { authStore } from '../services/auth'

interface ReviewsProps {
  locationId: string
  onReviewSubmitted?: () => void
}

function StarRating({
  value,
  onChange,
  size = 18,
  interactive = false,
}: {
  value: number
  onChange?: (v: number) => void
  size?: number
  interactive?: boolean
}) {
  const [hover, setHover] = useState(0)
  return (
    <div className="flex items-center gap-0.5">
      {[1, 2, 3, 4, 5].map((n) => {
        const filled = (hover || value) >= n
        return (
          <button
            key={n}
            type="button"
            disabled={!interactive}
            onClick={() => onChange?.(n)}
            onMouseEnter={() => interactive && setHover(n)}
            onMouseLeave={() => interactive && setHover(0)}
            className={interactive ? 'cursor-pointer' : 'cursor-default'}
            aria-label={`${n} star${n === 1 ? '' : 's'}`}
          >
            <Star
              size={size}
              className={filled ? 'text-yellow-500' : 'text-gray-300'}
              fill={filled ? 'currentColor' : 'none'}
            />
          </button>
        )
      })}
    </div>
  )
}

export default function Reviews({ locationId, onReviewSubmitted }: ReviewsProps) {
  const isAuthed = Boolean(authStore.getToken())

  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const [rating, setRating] = useState(0)
  const [comment, setComment] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [submitted, setSubmitted] = useState(false)

  const loadReviews = async () => {
    setLoading(true)
    setError(null)
    try {
      const items = await storageService.getLocationReviews(locationId)
      setReviews(items)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load reviews')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadReviews()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [locationId])

  const submit = async () => {
    if (rating < 1) {
      setSubmitError('Please select a rating')
      return
    }
    setSubmitting(true)
    setSubmitError(null)
    try {
      await storageService.createReview(locationId, { rating, comment: comment.trim() })
      setSubmitted(true)
      setRating(0)
      setComment('')
      await loadReviews()
      onReviewSubmitted?.()
    } catch (e) {
      setSubmitError(e instanceof Error ? e.message : 'Failed to submit review')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="flex items-center justify-between">
        <div className="text-lg font-semibold text-gray-900">Reviews</div>
        <div className="text-sm text-gray-600">{reviews.length} review{reviews.length === 1 ? '' : 's'}</div>
      </div>

      {/* Submit form */}
      {isAuthed ? (
        submitted ? (
          <div className="mt-4 rounded-md bg-green-50 px-4 py-3 text-sm text-green-700">
            Thanks for your review!
          </div>
        ) : (
          <div className="mt-4 rounded-md border border-gray-200 p-4">
            <div className="text-sm font-medium text-gray-900">Leave a review</div>
            <div className="mt-2">
              <StarRating value={rating} onChange={setRating} interactive size={24} />
            </div>
            <textarea
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder="Share your experience (optional)"
              rows={3}
              className="mt-3 w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-transparent focus:ring-2 focus:ring-blue-500"
            />
            {submitError && <div className="mt-2 text-sm text-red-600">{submitError}</div>}
            <button
              onClick={() => void submit()}
              disabled={submitting}
              className="mt-3 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              {submitting ? 'Submitting…' : 'Submit review'}
            </button>
          </div>
        )
      ) : (
        <div className="mt-4 rounded-md bg-gray-50 px-4 py-3 text-sm text-gray-600">
          Log in and book to leave a review.
        </div>
      )}

      {/* List */}
      <div className="mt-6 space-y-4">
        {loading && <div className="text-sm text-gray-600">Loading reviews…</div>}
        {error && <div className="text-sm text-red-600">{error}</div>}
        {!loading && !error && reviews.length === 0 && (
          <div className="text-sm text-gray-600">No reviews yet. Be the first to review!</div>
        )}
        {reviews.map((r) => (
          <div key={r.id} className="border-b border-gray-100 pb-4 last:border-b-0">
            <div className="flex items-center gap-3">
              <div className="flex h-9 w-9 items-center justify-center rounded-full bg-blue-600 text-sm font-semibold text-white">
                {(r.userName || 'U').charAt(0).toUpperCase()}
              </div>
              <div>
                <div className="text-sm font-medium text-gray-900">{r.userName || 'Traveler'}</div>
                <div className="text-xs text-gray-500">
                  {new Date(r.createdAt).toLocaleDateString(undefined, {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                  })}
                </div>
              </div>
              <div className="ml-auto">
                <StarRating value={r.rating} size={16} />
              </div>
            </div>
            {r.comment && <p className="mt-2 text-sm leading-relaxed text-gray-700">{r.comment}</p>}
          </div>
        ))}
      </div>
    </div>
  )
}
