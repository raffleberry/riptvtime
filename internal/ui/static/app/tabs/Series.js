import { useSeriesStore } from "../stores/series.js";
import { PAGE, theme } from "../utils.js";
import { onMounted, ref, storeToRefs, useRoute, watch } from "../vue.js";


const Series = {
    props: {

    },
    components: {

    },
    setup: (props) => {
        onMounted(() => {
        });

        const r = useRoute()
        const seriesStore = useSeriesStore()
        const { loading, seriesDetails, watchedEps } = storeToRefs(seriesStore)

        const { fetchSeries } = seriesStore

        const seriesId = ref(0)

        watch(() => r.params.id, (id) => {
            if (id) {
                fetchSeries(id)
            }
        }, { immediate: true })

        return {
            loading,
            sd: seriesDetails,
            watchedEps,
        }
    },
    template: `
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
                    <h3 class="card-title">{{ sd.Name }} <span class="text-muted">({{ sd.Year }}) </span></h3>
                    <p class="card-text">{{ sd.Overview }}</p>
                </div>
            </div>
            <div class="accordion" id="seasons">

                <div v-for="sn in sd.Seasons" class="accordion-item">
                    <h2 class="accordion-header">
                        <button class="accordion-button" type="button" data-bs-toggle="collapse" :data-bs-target="'#season' + sn.SeasonNumber">
                            {{ sn.Name }}
                        </button>
                    </h2>
                    <div :id="'season' + sn.SeasonNumber" class="accordion-collapse collapse show" data-bs-parent="#seasons">
                        <div class="accordion-body">
                            <div class="d-flex flex-row justify-content-between" v-for="ep in sn.Episodes">
                                <div>{{ ep.SeasonNumber }}x{{ ep.EpisodeNumber }} - {{ ep.Name }}</div>
                                <button class="btn">
                                    <i :class="(ep.SeasonNumber + 'x' + ep.EpisodeNumber) in watchedEps ? 'bi bi-check-circle-fill text-success' : 'bi bi-check-circle'"></i>
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

