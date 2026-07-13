import { ENDPOINT } from "../utils.js";
import { defineStore, ref, watch } from "../vue.js";

export const useSeriesStore = defineStore('series', () => {
    const loading = ref(false)
    const seriesDetails = ref({})
    const watchedEps = ref([])

    const fetchSeries = async (id) => {
        loading.value = true
        try {
            const response = await fetch(ENDPOINT.SERIES_GET(id));
            if (response.status === 200) {
                const result = await response.json();
                let we = result.EpsWatched
                delete result.EpsWatched
                seriesDetails.value = result

                if (we) {
                    let weArrStr = {}
                    for (let i = 0; i < we.length; i++) {
                        weArrStr[we[i].S + 'x' + we[i].E] = true
                    }
                    watchedEps.value = weArrStr
                }

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

    return {
        // data
        loading,
        seriesDetails,
        watchedEps,

        // actions
        fetchSeries,
    }

})