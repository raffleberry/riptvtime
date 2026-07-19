import { ENDPOINT, Ky } from "../utils.js";
import { computed, defineStore, ref, watch } from "../vue.js";

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

    const EpWatchCnt = computed(() => {
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

                // let weArrStr = {}
                // for (let i = 0; i < we.length; i++) {
                //     weArrStr[we[i].S + 'x' + we[i].E] = true
                // }
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

    const selectedEp = ref({})

    return {
        // data
        loading,
        seriesDetails,
        SnWatchedEps,
        EpWatchCnt,
        selectedEp,

        // actions
        fetchSeries,
    }

})