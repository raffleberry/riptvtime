import { apiGetImdbRating } from "../api.js"
import { defineStore, reactive, ref } from "../vue.js"

export const useImdb = defineStore("imdbStore", () => {
  const available = ref(true)
  const ratings = reactive({})

  const fetchRating = async (mId) => {
    try {
      const { data, err, statusCode } = await apiGetImdbRating(mId)
      if (err) {
        if (statusCode === 503) {
          console.warn("IMDB unavailable")
          available.value = false
          return
        }
        throw err
      }
      ratings[mId] = data
    } catch (error) {
      console.error("Error fetching feed data:", error)
    } finally {
    }
  }
  ;(async () => {
    await fetchRating(125988)
  })()

  const getRating = (mId) => {
    if (!available.value) return {}
    if (!ratings[mId]) fetchRating(mId)
    return ratings[mId]
  }

  return {
    available,
    getRating,
  }
})
