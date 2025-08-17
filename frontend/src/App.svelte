<script lang="ts">
  import Router, { link } from "svelte-spa-router";
  import Home from "./routes/Home.svelte"
  import Artists from "./routes/Artists.svelte"
  import Login from "./routes/Login.svelte"

  const routes = {
    "/": Home,
    "/artists": Artists,
    "/login": Login
  };

  console.log("App.svelte loaded");
</script>

<nav>
  <a href="/" use:link>Home</a>
  <a href="/artists" use:link>Artists</a>
</nav>

<main>
  <Router {routes} />
</main>

<style>
  nav {
    background: #f8f9fa;
    padding: 1rem;
    border-bottom: 1px solid #dee2e6;
  }

  nav a {
    margin-right: 1rem;
    text-decoration: none;
    color: #007bff;
  }

  nav a:hover {
    color: #0056b3;
  }

  main {
    padding: 2rem;
  }
</style>
