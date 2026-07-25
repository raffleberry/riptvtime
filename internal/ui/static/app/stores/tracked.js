import { apiAddSeries, apiRemSeries, apiSeriesTracked, apiSetStatus } from "../api.js";
import { notifyError } from "../components/Error.js";
import { ENDPOINT } from "../utils.js";
import { defineStore, ref, watch } from "../vue.js";

export const useTracked = defineStore('trackedSeries', () => {
    const loading = ref(false)
    const series = ref({})

    const refresh = async () => {
        loading.value = true
        try {
            const { data, err } = await apiSeriesTracked()
            if (err) {
                throw err
            }
            let ns = {}
            for (let i = 0; i < data.length; i++) {
                ns[data[i].MId] = data[i]
            }
            series.value = ns

        } catch (error) {
            console.error('Error fetching feed data:', error);
            notifyError(error)
        } finally {
            loading.value = false
        }
    }

    const addSeries = async (Id) => {

        try {

            const {data, err} = await apiAddSeries(Id)
            if (err) {
                throw err
            }
            series.value[Id] = data

        } catch (error) {
            console.error(error)
            return error
        } finally {
        }

        return null

    }
    const remSeries = async (Id) => {
    
        try {
            loading.value = true

            const {data, err} = await apiRemSeries(Id)
            if (err) {
                throw err
            }

            delete series.value[Id]

        } catch (error) {
            console.error(error)
            return error
        } finally {
            loading.value = false
        }

        return null

    }

    const changeStatus = async (Id, newStatus) => {

        try {
            const {data, err} = await apiSetStatus(Id, newStatus)
            if (err) {
                throw err
            }
            series.value[Id].TrackingStatus = newStatus
        } catch (error) {
            console.error('Error changing status:', error);
            notifyError(error)
        } finally {
        }

    }

    refresh()

    return {
        // data
        loading,
        series,

        // actions
        changeStatus,
        addSeries,
        remSeries,

    }

})