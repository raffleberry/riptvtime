import { apiEpUpNext, apiEpWatch } from "../../api.js"
import { MsgType, notify } from "../../components/Notify/Notify.js"
import { defineStore, ref } from "../../vue.js"

export const useFeedStore = defineStore("feed", () => {
  const loading = ref(false)
  const feed = ref([])

  const fetchFeed = async () => {
    loading.value = true
    try {
      const response = await fetch("/api/series/feed")
      if (response.status === 200) {
        const result = await response.json()
        feed.value = result
      } else {
        const msg = `${response.status} - ${await response.text()}`
        throw new Error(msg)
      }
    } catch (error) {
      console.error("Error fetching feed data:", error)
    } finally {
      loading.value = false
    }
  }

  const epMarkAndGetUpNext = async (mId, sNo, epNo) => {
    try {
      let err1 = await apiEpWatch(mId, [{ S: sNo, E: epNo }])
      if (err1) {
        throw err1
      }

      let { data, err } = await apiEpUpNext(mId)
      if (err) {
        throw err
      }

      // process
      if (data.S === -1 || data.E === -1) {
        feed.value = feed.value.filter((f) => f.MId !== mId)
      } else {
        feed.value = feed.value.map((f) => {
          if (f.MId === mId) {
            f.UpNextS = data.S
            f.UpNextE = data.E
            f.EpisodesWatched = f.EpisodesWatched + 1
          }
          return f
        })
      }
    } catch (error) {
      console.error(mId, sNo, epNo, "Error marking ep watched:", error)
      notify(MsgType.Error, "Feed", error)
    } finally {
    }
  }

  return {
    // data
    loading,
    feed,

    // actions
    fetchFeed,
    epMarkAndGetUpNext,
  }
})
