import { apiAddSeries, apiRemSeries, apiSetStatus } from "../../api.js";
import { ENDPOINT, TvStatus } from "../../utils.js";
import { defineStore, ref, storeToRefs } from "../../vue.js";
import { useSearchStore } from "./searchStore.js";

export const useSeriesOpts = defineStore('SeriesOptsStore', () => {
    const loading = ref(false)

    const selected = ref({
        Id: null,
        Name: null,
        Year: null,
        Status: 0,
    })


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

            const {data, err} = await apiAddSeries(selected.value.Id)
            if (err) {
                throw err
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

            const {data, err} = await apiRemSeries(selected.value.Id)
            if (err) {
                throw err
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



