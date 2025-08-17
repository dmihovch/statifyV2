<script lang="ts">
  import { onMount } from 'svelte';
  import { link } from 'svelte-spa-router';
  
  let isLoggedIn = false;
  let spotifyCookie: string | null = null;
  
  onMount(() => {
    checkLoginStatus();
  });
  
  function checkLoginStatus() {
    const cookies = document.cookie.split(';');
    spotifyCookie = cookies.find(cookie => 
      cookie.trim().startsWith('spotify_user_id=')
    ) || null;
    
    if (spotifyCookie) {
      isLoggedIn = true;
    }
  }
  
  function logout() {
    document.cookie = 'spotify_user_id=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    isLoggedIn = false;
    spotifyCookie = null;
    window.location.reload(); // Refresh to update the UI
  }
</script>

<div>
  <h1>Welcome to Statify</h1>
  
  {#if isLoggedIn}
    <div class="logged-in">
      <p>You are logged in! 🎉</p>
      <button on:click={logout} class="logout-btn">Logout</button>
    </div>
  {:else}
    <div class="not-logged-in">
      <p>Please log in to see your Spotify statistics</p>
      <a href="/login" use:link class="login-link">Login with Spotify</a>
    </div>
  {/if}
</div>

<style>
  div {
    padding: 1rem;
  }
  
  h1 {
    color: #333;
    margin-bottom: 1rem;
  }
  
  .logged-in {
    background: #d4edda;
    border: 1px solid #c3e6cb;
    border-radius: 8px;
    padding: 20px;
    margin: 20px 0;
  }
  
  .not-logged-in {
    background: #f8f9fa;
    border: 1px solid #dee2e6;
    border-radius: 8px;
    padding: 20px;
    margin: 20px 0;
  }
  
  .logout-btn {
    background: #dc3545;
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .logout-btn:hover {
    background: #c82333;
  }
  
  .login-link {
    display: inline-block;
    background: #1DB954;
    color: white;
    text-decoration: none;
    padding: 10px 20px;
    border-radius: 25px;
    margin-top: 10px;
  }
  
  .login-link:hover {
    background: #1ed760;
  }
</style>
