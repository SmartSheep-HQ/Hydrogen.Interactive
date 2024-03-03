import { createRouter, createWebHistory } from "vue-router"
import MasterLayout from "@/layouts/master.vue"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: MasterLayout,
      children: [
        {
          path: "/",
          name: "explore",
          component: () => import("@/views/explore.vue")
        },

        {
          path: "/p/:postType/:alias",
          name: "posts.details",
          component: () => import("@/views/posts/details.vue")
        }
      ]
    }
  ]
})

export default router
