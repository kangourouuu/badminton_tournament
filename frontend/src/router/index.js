import { createRouter, createWebHistory } from "vue-router";
import TournamentMap from "../views/TournamentMap.vue";
import LoginView from "../views/LoginView.vue";
import AdminDashboard from "../views/AdminDashboard.vue";
import RulesView from "../views/RulesView.vue";
import ContactView from "../views/ContactView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: TournamentMap,
    },
    {
      path: "/login",
      name: "login",
      component: LoginView,
    },
    {
      path: "/rules",
      name: "rules",
      component: RulesView,
    },
    {
      path: "/contact",
      name: "contact",
      component: ContactView,
    },
    {
      path: "/admin",
      name: "admin",
      component: AdminDashboard,
      meta: { requiresAuth: true },
    },
  ],
});

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem("token");
    if (!token) {
      next({ name: "login" });
    } else {
      next();
    }
  } else {
    next();
  }
});

export default router;
