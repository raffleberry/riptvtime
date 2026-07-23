import { SeriesOptions } from "../../overlays/SeriesOptions.js";
import { useConfirm } from "../../stores/confirm.js";
import { Ky, PAGE, theme } from "../../utils.js";
import { computed, onMounted, ref, storeToRefs, useRoute, watch } from "../../vue.js";
import { useSeriesStore } from "./seriesStore.js";


const Series = {
    props: {

    },
    components: {
        SeriesOptions
    },
    setup: (props) => {
        onMounted(() => {
        });

        const r = useRoute()
        const seriesStore = useSeriesStore()
        const { loading, seriesDetails, SnWatchedEps, EpWatchCnt } = storeToRefs(seriesStore)

        const confirmStore = useConfirm()
        const { openDialog } = confirmStore

        const { fetchSeries } = seriesStore

        const seriesId = ref(0)

        watch(() => r.params.id, (id) => {
            if (id) {
                fetchSeries(id)
            }
        }, { immediate: true })

        return {
            Ky,
            loading,
            sd: seriesDetails,
            SnWatchedEps,
            EpWatchCnt,
            openDialog,
        }
    },
    template: `
    <SeriesOptions></SeriesOptions>
    <div class="container-fluid">
        <div v-if="loading" class="d-flex justify-content-center align-items-center"
            style="min-height: 50vh;">
            <div class="spinner-border" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
        <div v-else>
            <div class="card">
                <div class="card-body">
                    <div class="d-flex justify-content-between">
                        <h3 class="card-title">{{ sd.Name }} <span class="text-muted">({{ sd.Year }}) </span></h3>
                        <button @click="openSeriesOptions" 
                            data-bs-toggle="offcanvas" data-bs-target="#seriesOptions"
                            type="button" class="btn p-2 d-inline-flex align-items-center justify-content-center">
                                <i :class="statusCss.icon"></i>
                        </button>
                    </div>
                    <p class="card-text">{{ sd.Overview }}</p>
                </div>
            </div>
            <div class="accordion" id="seasons">

                <div v-for="sn in sd.Seasons" class="accordion-item">
                    <h2 class="accordion-header">
                        <button class="accordion-button" type="button" data-bs-toggle="collapse" :data-bs-target="'#season' + sn.SeasonNumber">
                            <div class="d-flex w-100 justify-content-between pe-5"><span>{{ sn.Name }}</span> <span>{{SnWatchedEps[sn.SeasonNumber]?.length}}/{{sn.EpisodeCount}}</span></div>
                        </button>
                    </h2>
                    <div :id="'season' + sn.SeasonNumber" class="accordion-collapse collapse" data-bs-parent="#seasons">
                        <div class="accordion-body">
                            <div class="d-flex flex-row justify-content-between" v-for="ep in sn.Episodes">
                                <div>{{ ep.SeasonNumber }}x{{ ep.EpisodeNumber }} - {{ ep.Name }}</div>
                                <button class="btn"
                                    @click="() => { openDialog(Ky(ep.SeasonNumber, ep.EpisodeNumber)) }"
                                >
                                    <i :class="Ky(ep.SeasonNumber, ep.EpisodeNumber) in EpWatchCnt ? 'bi bi-check-circle-fill text-success' : 'bi bi-check-circle'"></i>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    </div>
    `
}
export { Series };

