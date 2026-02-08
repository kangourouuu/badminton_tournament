import { createRouter, createWebHashHistory } from "vue-router";
import PublicView from "../components/PublicView.vue";
import AdminDashboard from "../components/AdminDashboard.vue";
import LoginView from "../views/LoginView.vue";

const routes = [
  {
    path: "/",
    name: "public",
    component: PublicView,
  },
  {
    path: "/admin/login",
    name: "login",
    component: LoginView,
  },
  {
    path: "/admin",
    name: "admin",
    component: AdminDashboard,
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem("token");
  if (to.matched.some((record) => record.meta.requiresAuth)) {
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
