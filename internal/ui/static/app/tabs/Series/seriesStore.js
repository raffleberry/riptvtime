import { apiEpUnWatch, apiEpWatch } from "../../api.js";
import { ENDPOINT, Ky } from "../../utils.js";
import { computed, defineStore, ref, watch } from "../../vue.js";

export const useSeriesStore = defineStore('series', () => {
    const loading = ref(false)
    const seriesDetails = ref({})
    const watchedEps = ref([])
    const SnWatchedEps = computed(() => {
        let rv = {}
        for (let i = 0; i <= seriesDetails.value.NumberOfSeasons; i++) {
            rv[i] = []
        }
        for (const ep of watchedEps.value) {
            if (rv[ep.S].includes(ep.E)) {
                continue
            }
            rv[ep.S].push(ep.E)
        }
        return rv
    })

    const epWatchCnt = computed(() => {
        let rv = {}
        for (const ep of watchedEps.value) {
            if (!rv[Ky(ep.S, ep.E)]) {
                rv[Ky(ep.S, ep.E)] = 0
            }
            rv[Ky(ep.S, ep.E)] += 1
        }
        return rv
    })

    const fetchSeries = async (id) => {
        loading.value = true
        try {
            const response = await fetch(ENDPOINT.SERIES_GET(id));
            if (response.status === 200) {
                const result = await response.json();
                let we = result.EpsWatched
                delete result.EpsWatched
                seriesDetails.value = result

                watchedEps.value = we || []

            } else {
                const msg =`${response.status} - ${await response.text()}`
                throw new Error(msg)
            }
        } catch (error) {
            console.error('Error fetching series data:', error);
        } finally {
            loading.value = false
        }
    }

    const epMarkWatched = async (mId, sNo, epNo) => {
        try {
            const err = await apiEpWatch(mId, sNo, epNo);
            if (err) {
                throw err
            }

            watchedEps.value.push({ S: sNo, E: epNo })

        } catch (error) {
            console.error('Error fetching series data:', error);
        } finally {
        }
    }

    const epUnMarkWatched = async (mId, sNo, epNo) => {
        try {
            const err = await apiEpUnWatch(mId, sNo, epNo);
            if (err) {
                throw err
            }

            const idx = watchedEps.value.findIndex(ep => ep.S === sNo && ep.E === epNo)
            if (idx !== -1) {
                watchedEps.value.splice(idx, 1)
            } else {
                console.error(EpsWatched)
                throw new Error('Episode not found in watched list')
            }

        } catch (error) {
            console.error('Error fetching series data:', error);
        } finally {
        }
    }

    const selectedEp = ref({})

    return {
        // data
        loading,
        seriesDetails,
        SnWatchedEps,
        epWatchCnt,
        selectedEp,

        // actions
        fetchSeries,
        epMarkWatched,
        epUnMarkWatched,
    }

})