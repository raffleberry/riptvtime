import { ky, TvStatus } from "../../utils.js"
import { computed, onMounted, ref, storeToRefs, watch } from "../../vue.js"
import { notifyError } from "../../components/Error.js"
import { useTracked } from "../../stores/tracked.js"
import { useSeriesStore } from "./seriesStore.js"

export const EpisodeOpts = {
    props: {
        mid: Number,
        ep: Object,
    },

    setup(props) {
        const loading = ref(false)

        var bSelf = null

        onMounted(() => {
            const el = document.getElementById('episodeOpts')
            bSelf = bootstrap.Offcanvas.getOrCreateInstance(el)

            // el.addEventListener('show.bs.offcanvas', () => {
            // })

            // el.addEventListener('hidden.bs.offcanvas', () => {
            // })
        })
        const store = useSeriesStore()

        const { epMarkWatched, epUnMarkWatched } = store

        const { epWatchCnt } = storeToRefs(store)

        const handleIncr = async () => {
            try {
                loading.value = true
                await epMarkWatched(props.mid, props.ep.SeasonNumber, props.ep.EpisodeNumber)
            } catch (e) {
            } finally {
                bSelf.hide()
                loading.value = false
            }
        }        
        const handleDecr = async () => {
            try {
                loading.value = true
                await epUnMarkWatched(props.mid, props.ep.SeasonNumber, props.ep.EpisodeNumber)
            } catch (e) {
            } finally {
                bSelf.hide()
                loading.value = false
            }
        }

        const cnt = computed(() => epWatchCnt.value[ky(props.ep.SeasonNumber, props.ep.EpisodeNumber)] ?? 0 )

        const getIncrTxt = (c) => {
            if (c == 0) {
                return "Watch"
            } else {
                return "Re-Watch"
            }
        }

        const getDecrTxt = (c) => {
            if (c == 1) {
                return "Remove Watched"
            } else {
                return "Remove Re-Watched"
            }
        }

        return {
            loading,
            handleIncr,
            handleDecr,
            getIncrTxt,
            getDecrTxt,
            cnt
        }

    },

    template: `
    <div class="offcanvas offcanvas-end" tabindex="-1" id="episodeOpts">
        <div class="offcanvas-header border-bottom">
            <h5 class="offcanvas-title"><span class="text-muted">{{ep.Name}}</span></h5>
            <button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close"></button>
        </div>
        <div class="offcanvas-body p-0">
            <p class="px-3 mt-2">Season: {{ep.SeasonNumber}} Episode: {{ep.EpisodeNumber}}</p>
            <p class="px-3 fst-italic">{{ep.Runtime}} minutes</p>
            <p class="px-3">{{ep.Overview}}</p>
            <div v-if="loading" class="list-group list-group-flush" >
            <div class="list-group-item list-group-item-action px-4 py-3 d-flex align-items-center justify-content-center border-0">
            <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
            </div>
            </div>
            </div>
            <div v-else class="list-group list-group-flush">
                <div class="px-3" class="list-group-item"> <span v-if="cnt > 1">{{cnt}}x</span> <i class="bi bi-check-circle-fill text-success"></i></div>
                <button :disabled="loading" class="list-group-item list-group-item-action btn btn-primary px-4 py-3 d-flex align-items-center border-0 text-primary"
                    @click="handleIncr">
                    <i class="bi bi-plus-circle me-2"></i> {{ getIncrTxt(cnt) }}
                </button>
                <button :disabled="loading" class="list-group-item list-group-item-action btn btn-primary px-4 py-3 d-flex align-items-center border-0 text-danger"
                    @click="handleDecr"> 
                    <i class="bi bi-dash-circle me-2"></i> {{ getDecrTxt(cnt) }}
                </button>
            </div>
        </div>
    </div>
`
}
