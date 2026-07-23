import { useSeriesOpts } from "../tabs/Search/optsStore.js"
import { TvStatus } from "../utils.js"
import { computed, onMounted, ref, storeToRefs } from "../vue.js"

const SeriesOptions = {

    setup() {

        const store = useSeriesOpts()

        const { selected, loading } = storeToRefs(store)

        const { changeStatus, addSeries, remSeries  } = store

        const onDetails = () => {
            console.log("Details")
        }


        const statusActionTxt = computed(() => {
            if (selected.value.Status === TvStatus.Watching) {
                return 'Stop watching'
            } else if (selected.value.Status === TvStatus.Stopped) {
                return 'Resume watching'
            }
        })

        const addRemActionTxt = computed(() => {
            if (selected.value.Status === TvStatus.NotWatching) {
                return 'Add to list'
            } else {
                return 'Remove from list'
            }
        })

        var bSelf = null

       
        onMounted(() => {
            const el = document.getElementById('seriesOptions')
            bSelf = new bootstrap.Offcanvas(el)



            // el.addEventListener('show.bs.offcanvas', () => {
            // })

            el.addEventListener('hidden.bs.offcanvas', () => {
                selected.value = {}
            })
        })

        const handleStatusChange = async () => {
            await changeStatus()
            bSelf.hide()
        }

        const handleAddRemSeries = async () => {
            if (selected.value.Status === TvStatus.NotWatching) {
                const err = await addSeries()
                if (err) {
                    //TODO: notify error 
                } else {
                    selected.value.Status = TvStatus.Watching
                }
            } else {
                const err = await remSeries()
                if (err) {
                    //TODO: notify error
                } else {
                    selected.value.Status = TvStatus.NotWatching
                }
            }

            bSelf.hide()
        }
        // TvStatus.

        return {
            selected,
            loading,
            TvStatus,

            statusActionTxt,
            addRemActionTxt,

            handleStatusChange,
            handleAddRemSeries,

            onDetails,

        }

    },

    template: `
    <div class="offcanvas offcanvas-end" tabindex="-1" id="seriesOptions">
        <div class="offcanvas-header border-bottom">
            <h5 class="offcanvas-title">{{selected.Name}} <span class="text-muted">({{selected.Year}})</span></h5>
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
                <button class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-primary">
                    Show details
                </button>
                <button class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-primary"
                    v-if="[TvStatus.Watching, TvStatus.Stopped].includes(selected.Status)" class="btn btn-primary" @click="handleStatusChange">
                    {{ statusActionTxt }}
                </button>
                <button class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center border-0 text-danger"
                    class="btn btn-primary" @click="handleAddRemSeries">
                    {{ addRemActionTxt }}
                </button>
            </div>
        </div>
    </div>
`
}

export { SeriesOptions }