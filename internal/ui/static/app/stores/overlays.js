import { ENDPOINT, TvStatus } from "../utils.js";
import { defineStore, ref, storeToRefs } from "../vue.js";
import { useSearchStore } from "./search.js";

export const useSeriesOpts = defineStore('SeriesOptsStore', () => {
    const loading = ref(false)

    const selected = ref({
        Id: null,
        Name: null,
        Year: null,
        Status: 0,
    })


    const apiSetStatus = async (Id, newStatus) => {
        let url = `${ENDPOINT.SERIES_STATUS(Id)}`
        try {
            const response = await fetch(url, {
                method: 'PUT',
                headers: { 'Content-Type' : 'application/json' },
                body: JSON.stringify({ status: newStatus })
            });

            if (!response.ok) {
                throw new Error(`Error, bad response from server: ${response.status} - ${response.statusText}`)
            }
            
            return {
                data: "OK",
                err: null
            }
        } catch (error) {
            console.error('error updating status:', error);
            return {
                data: null,
                err: error
            }
        }
    }

     const changeStatus = async () => {
        console.log("Clicked", loading.value, selected.value.Id)

        if (!selected.value.Id || loading.value) return
        console.log(!selected.value.Id || loading.value, selected.value)
        try {
            loading.value = true

            const { results, pageCur } =  storeToRefs(useSearchStore())

            const findAndUpdate = (id, newStatus) => {
                for (let i = 0; i < results.value[pageCur.value].length; i++) {
                    if (results.value[pageCur.value].Id === id) {
                        console.log("updating")
                        results.value[pageCur.value].Status = newStatus
                    }
                }
            }


            if (selected.value.Status === TvStatus.Watching) {

                const {data, err} = await apiSetStatus(selected.value.Id, TvStatus.Stopped)
                if (err) {
                    throw err
                }

                findAndUpdate(selected.value.Id, TvStatus.Stopped)

                //TODO:update feed store too... 

                // set selected status to the newer one
                selected.value.Status = TvStatus.Stopped
            } else if (selected.value.Status === TvStatus.Stopped) {

                const {data, err} = await apiSetStatus(selected.value.Id, TvStatus.Watching)
                if (err) {
                    throw err
                }

                findAndUpdate(selected.value.Id, TvStatus.Watching)

                //TODO:update feed store too... 

                // set selected status to the newer one
                selected.value.Status = TvStatus.Watching

            }
    
        } catch (error) {
            console.error(error)
        } finally {
            loading.value = false
        }

 
    }

    const addSeries = async () => {

        if (!selected.value.Id || loading.value) return

        try {
            loading.value = true

            const res = await fetch(ENDPOINT.SERIES_ADD(), {
                method: 'POST',
                headers: { 'Content-Type' : 'application/json' },
                body: JSON.stringify({
                    m_id: selected.value.Id,
                    m_name: selected.value.MName
                })
            })

            if (res.status !== 200) {
                console.error(await res.text())
                throw new Error(`Error, bad response from server: ${res.status} - ${res.statusText}`)
            }


            const { results, pageCur } =  storeToRefs(useSearchStore())

            const findAndUpdate = (id, newStatus) => {
                for (let i = 0; i < results.value[pageCur.value].length; i++) {
                    if (results.value[pageCur.value].Id === id) {
                        console.log("updating")
                        results.value[pageCur.value].Status = newStatus
                    }
                }
            }

            findAndUpdate(selected.value.Id, TvStatus.Watching)
            

        } catch (error) {
            console.error(error)
            return error
        } finally {
            loading.value = false
        }

        return null

    }

    const remSeries = async () => {

        if (!selected.value.Id || loading.value) return

        try {
            loading.value = true

            // remove from feed and me data

            const res = await fetch(ENDPOINT.SERIES_REM(selected.value.Id), {
                method: 'DELETE',
            })

            if (res.status !== 200 && res.status !== 204 && res.status !== 404) {
                console.error(await res.text())
                throw new Error(`Error, bad response from server: ${res.status} - ${res.statusText}`)
            }

            const { results, pageCur } =  storeToRefs(useSearchStore())

            const findAndUpdate = (id, newStatus) => {
                for (let i = 0; i < results.value[pageCur.value].length; i++) {
                    if (results.value[pageCur.value].Id === id) {
                        console.log("updating")
                        results.value[pageCur.value].Status = newStatus
                    }
                }
            }

            findAndUpdate(selected.value.Id, TvStatus.NotWatching)
            

        } catch (error) {
            console.error(error)
            return error
        } finally {
            loading.value = false
        }

        return null

    }

    return {
        // data
        selected,
        loading,

        // actions
        changeStatus,
        remSeries,
        addSeries

    }

})



