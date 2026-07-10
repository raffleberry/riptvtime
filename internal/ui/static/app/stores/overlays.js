import { ENDPOINT } from "../utils.js";
import { defineStore, ref } from "../vue.js";

export const useSeriesOpts = defineStore('SeriesOptsStore', () => {
    const loading = ref(false)

    const selected = ref({
        MId: null,
        Name: null,
        Year: null,
    })

    let status = -1

    const setStatus = async (mId, newStatus) => {
        let url = `${ENDPOINT.SERIES_STATUS(mId)}`
        try {
            const response = await fetch(url, {
                method: 'PUT',
                headers: { 'Content-Type' : 'application/json' },
                body: JSON.stringify({ status: newStatus })
            });
            const result = await response.json();
            return {
                data: result,
                err: null
            }
        } catch (error) {
            console.error('error updating status:', error);
            return {
                data: result,
                err: error
            }
        }
    }

     const onStatus = async () => {
        console.log("Clicked")
        return
        
        if (!selected.value.mId || loading.value) return

        if (status === TvStatus.Watching) {

            // set selected status to the newer one

            // do the appropriate api call to make the changes

            // update the search results data

            status = TvStatus.NotWatching
        } else if (status === TvStatus.NotWatching) {
            status = TvStatus.Watching
        }
    
        try {
            loading.value = true
            const { data, err } = await searchTv(newStatus, 1)
            if (err) {
                throw err
            }

            results.value = {}

        } catch (error) {
            console.error(error)
        } finally {
            loading.value = false
        }
    
    }

    return {
        // data
        selected,
        loading,

        // actions
        onStatus,

    }

})



