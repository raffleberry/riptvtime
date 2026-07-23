import { Ky } from "../../utils.js"
import { TvStatus } from "../utils.js"
import { computed, onMounted, ref, storeToRefs } from "../vue.js"
import { useSeriesStore } from "./seriesStore.js"

export const EpisodeOptions = {

    setup() {

        const store = useSeriesStore()

        const { selectedEp } = storeToRefs(store)

        const { markWatch, markUnwatch } = store


        var elBs = null

       
        onMounted(() => {
            const el = document.getElementById('episodeOptions')
            elBs = new bootstrap.Offcanvas(el)

            el.addEventListener('hidden.bs.offcanvas', () => {
                selectedEp.value = {}
            })
        })

        const handleMarkWatch = async () => {
            await markWatch()
            elBs.hide()
        }

        
        const handleMarkUnwatch = async () => {
            await markUnwatch()
            elBs.hide()
        }


        return {
            selectedEp,
            loading,
            TvStatus,

            statusActionTxt,
            addRemActionTxt,

            handleStatusChange,
            handleAddRemSeries,

            onDetails,

            Ky

        }

    },

    template: `
    <div class="offcanvas offcanvas-end" tabindex="-1" id="episodeOptions">
        <div class="offcanvas-header border-bottom">
            <h5 class="offcanvas-title"><span class="badge bg-primary">({{Ky(selectedEp.SeasonNumber, selectedEp.EpisodeNumber)}})</span> {{selectedEp.Name}}</h5>
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
                    class="btn btn-primary" @click="handleStatusChange">
                    {{ if  }}
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

