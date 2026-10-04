import { imgPosterUrl } from "../utils.js"

const DiscoverRail = {
  props: {
    title: String,
    listKey: String,
    items: Array,
    loading: Boolean,
  },
  setup() {
    return {
      imgPosterUrl,
    }
  },
  template: /* HTML */ `
    <div class="mb-4">
      <div class="d-flex justify-content-between align-items-center px-3 mb-2">
        <h4 class="m-0">{{ title }}</h4>
        <router-link
          :to="{ path: '/discover', query: { view: listKey, p: 1 } }"
          class="btn btn-sm btn-outline-primary"
        >
          View All
        </router-link>
      </div>
      <div
        v-if="loading"
        class="d-flex justify-content-center align-items-center"
        style="min-height: 20vh;"
      >
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      <div v-else class="d-flex flex-row overflow-auto gap-3 px-3 pb-2">
        <router-link
          v-for="tv in items"
          :key="tv.Id"
          :to="'/series/' + tv.Id"
          class="text-decoration-none flex-shrink-0"
          style="width: 150px;"
        >
          <img
            :src="imgPosterUrl(tv.Image)"
            class="img-fluid rounded"
            style="width: 150px; height: 225px; object-fit: cover;"
            loading="lazy"
            alt="poster"
          />
          <div class="mt-1 text-truncate text-body">
            {{ tv.Name }}
            <span class="text-muted">({{ tv.Year }})</span>
          </div>
        </router-link>
      </div>
    </div>
  `,
}

export { DiscoverRail }
