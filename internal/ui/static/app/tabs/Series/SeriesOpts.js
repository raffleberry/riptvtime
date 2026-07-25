import { TvStatus } from "../../utils.js"
import { computed, onMounted, ref, storeToRefs, watch } from "../../vue.js"
import { notifyError } from "../../components/Error.js"
import { useTracked } from "../../stores/tracked.js"

export const SeriesOpts = {
    props: {
        mid: Number,
        name: String,
        year: Number,
    },
    setup(props) {
        const loading = ref(false)

        const trackedStore = useTracked()

        const { series } = storeToRefs(trackedStore)
        const { changeStatus, remSeries, addSeries } = trackedStore

        const status = computed(() => {
            let ob = series.value?.[props.mid]
            if (ob) {
                return ob.TrackingStatus
            }
            return TvStatus.NotWatching
        })

        const getStatusBtnTxt = (s) => (s === TvStatus.Watching) ? 'Stop watching' : 'Resume watching'
        const getAddRemBtnTxt = (s) => (s === TvStatus.NotWatching) ? 'Add Series' : 'Remove Series'
        

        var bSelf = null


        onMounted(() => {
            const el = document.getElementById('seriesOpts')
            bSelf = new bootstrap.Offcanvas(el)

            // el.addEventListener('show.bs.offcanvas', () => {
            // })

            // el.addEventListener('hidden.bs.offcanvas', () => {
            // })
        })

        const handleStatusChange = async () => {
            try {
                loading.value = true

                if (status.value === TvStatus.Watching) {
                    await changeStatus(props.mid, TvStatus.Stopped)
                } else if (status.value === TvStatus.Stopped) {
                    await changeStatus(props.mid, TvStatus.Watching)
                }
                bSelf.hide()
            } catch (error) {
                notifyError(error)
            } finally {
                loading.value = false
            }
        }

        const handleAddRemSeries = async () => {
            try {
                loading.value = true
                if (status.value === TvStatus.NotWatching) {
                    const err = await addSeries(props.mid)
                    if (err) {
                        throw (err)
                    }
                } else {
                    const err = await remSeries(props.mid)
                    if (err) {
                        throw (err)
                    }
                }

                bSelf.hide()
            } catch (error) {
                console.error(error)
                notifyError(error)
            } finally {
                loading.value = false
            }

        }

        return {
            loading,
            TvStatus,
            status,

            getStatusBtnTxt,
            getAddRemBtnTxt,

            handleStatusChange,
            handleAddRemSeries,


        }

    },

    template: `
    <div class="offcanvas offcanvas-end" tabindex="-1" id="seriesOpts">
        <div class="offcanvas-header border-bottom">
            <h5 class="offcanvas-title">{{name}} <span class="text-muted">({{year}})</span></h5>
            <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close"></button>
        </div>
        <div class="offcanvas-body p-0">
            <div v-if="loading" class="list-group list-group-flush" >
                <div class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center justify-content-center border-0">
                    <div class="spinner-border" role="status">
                        <span class="visually-hidden">Loading...</span>
                    </div>
                </div>
            </div>
            <div v-else class="list-group list-group-flush">
                <button class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-primary"
                    v-if="[TvStatus.Watching, TvStatus.Stopped].includes(status)" class="btn btn-primary" @click="handleStatusChange">
                    {{ getStatusBtnTxt(status) }}
                </button>
                <button class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-danger"
                    class="btn btn-primary" @click="handleAddRemSeries">
                    {{ getAddRemBtnTxt(status) }}
                </button>
            </div>
        </div>
    </div>
`
}
