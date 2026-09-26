import { apiEpUnWatch, apiEpWatch, apiGetSeriesDetails } from "../../api.js"
import { MsgType, notify } from "../../components/Notify/Notify.js"
import { useTracked } from "../../stores/tracked.js"
import { TvStatus } from "../../utils.js"
import { computed, defineStore, ref, storeToRefs } from "../../vue.js"

export const useSeriesStore = defineStore("series", () => {
  const loading = ref(false)
  const sd = ref({})
  const watchedEps = ref([])
  const watchedSpls = ref([])

  window.we = watchedEps
  const SnWatchedEps = computed(() => {
    let rv = {}
    for (let i = 0; i <= sd.value.NumberOfSeasons; i++) {
      rv[i] = []
    }
    for (const ep of watchedEps.value) {
      if (rv[ep.S].includes(ep.E)) {
        continue
      }
      rv[ep.S].push(ep.E)
    }
    for (const ep of watchedSpls.value) {
      if (rv[ep.S].includes(ep.E)) {
        continue
      }
      rv[ep.S].push(ep.E)
    }
    return rv
  })
  const tstore = useTracked()
  const { series: tSeries } = storeToRefs(tstore)
  const { addSeries: tAddSeries, remSeries: tRemSeries } = tstore

  const updateTrackingStore = (mId) => {
    if (!tSeries.value[mId]) return
    if (tSeries.value[mId].TrackingStatus === TvStatus.Stopped) {
      return
    }
    if (watchedEps.value.length === sd.value.EpisodesAired) {
      if (sd.value.InProduction) {
        tSeries.value[mId].TrackingStatus = TvStatus.UpToDate
      } else {
        tSeries.value[mId].TrackingStatus = TvStatus.Completed
      }
    } else {
      tSeries.value[mId].TrackingStatus = TvStatus.Watching
    }
  }

  const fetchSeries = async (id) => {
    loading.value = true
    try {
      const { data, err } = await apiGetSeriesDetails(id)
      if (err) {
        throw err
      }

      let weps = []
      let wspls = []
      for (const ep of data.EpsWatched) {
        if (ep.S === 0) {
          wspls.push(ep)
        } else {
          weps.push(ep)
        }
      }

      delete data.EpsWatched

      sd.value = data

      watchedEps.value = weps
      watchedSpls.value = wspls
    } catch (error) {
      console.error("Error getting series data:", error)
      notify(MsgType.Error, "Series", error)
    } finally {
      loading.value = false
    }
  }

  const epMarkWatched = async (mId, eps) => {
    try {
      if (!tSeries.value?.[mId]) {
        const err = await tAddSeries(mId)
        if (err) {
          throw err
        }
      }

      const err = await apiEpWatch(mId, eps)
      if (err) {
        throw err
      }
      for (const ep of eps) {
        let watched = watchedEps
        if (ep.S === 0) {
          watched = watchedSpls
        }

        let idx = watched.value.findIndex((sep) => sep.S === ep.S && sep.E === ep.E)
        if (idx === -1) {
          watched.value.push({
            S: ep.S,
            E: ep.E,
            Cnt: 1,
          })
        } else {
          watched.value[idx].Cnt += 1
        }
      }
      updateTrackingStore(mId)
    } catch (error) {
      console.error("Error fetching series data:", error)
      notify(MsgType.Error, "Series", error)
    } finally {
    }
  }

  const epUnMarkWatched = async (mId, sNo, epNo) => {
    try {
      const err = await apiEpUnWatch(mId, sNo, epNo)
      if (err) {
        throw err
      }
      let watched = watchedEps
      if (sNo === 0) {
        watched = watchedSpls
      }

      const idx = watched.value.findIndex((ep) => ep.S === sNo && ep.E === epNo)
      if (idx !== -1) {
        if (watched.value[idx].Cnt > 1) {
          watched.value[idx].Cnt -= 1
        } else {
          watched.value.splice(idx, 1)
          updateTrackingStore(mId)
        }
      } else {
        console.error(watched)
        throw new Error("Episode not found in watched list")
      }
    } catch (error) {
      console.error("Error fetching series data:", error)
      notify(MsgType.Error, "Series", error)
    } finally {
    }
  }

  const addSeries = async (mId) => {
    try {
      const err = await tAddSeries(mId)
      if (err) {
        throw err
      }
      updateTrackingStore(mId)
    } catch (error) {
      console.error(error)
      notify(MsgType.Error, "Series", error)
    } finally {
    }
  }

  const remSeries = async (mId) => {
    try {
      const err = await tRemSeries(mId)
      if (err) {
        throw err
      }
    } catch (error) {
      console.error(error)
      notify(MsgType.Error, "Series", err)
    } finally {
    }
  }

  const selectedEp = ref({})

  return {
    // data
    loading,
    sd,
    SnWatchedEps,
    selectedEp,
    watchedEps,
    watchedSpls,

    // actions
    fetchSeries,

    epMarkWatched,
    epUnMarkWatched,

    addSeries,
    remSeries,
  }
})
