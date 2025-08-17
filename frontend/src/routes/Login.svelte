<script lang="ts">
    function loginWithSpotify() {
        const clientId = '402291dff4ef472f9ae716691a4d77b7';
        const redirectUri = 'http://127.0.0.1:3001/login/callback';
        
        // These are the permissions we want from the user
        const scope = 'user-read-private user-read-email user-top-read user-read-recently-played';
        
        // Build the Spotify authorization URL with proper encoding
        const authUrl = `https://accounts.spotify.com/authorize?client_id=${clientId}&response_type=code&redirect_uri=${encodeURIComponent(redirectUri)}&scope=${encodeURIComponent(scope)}&show_dialog=true`;
        
        // Redirect the user to Spotify
        window.location.href = authUrl;
    }
</script>

<div class="login-container">
    <div class="login-card">
        <h1>Welcome to Statify</h1>
        <p>Connect your Spotify account to see your listening statistics</p>

        <button class="spotify-login-btn" on:click={loginWithSpotify}>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                <path
                    d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z"
                />
            </svg>
            Login with Spotify
        </button>

        <div class="features">
            <h3>What you'll get:</h3>
            <ul>
                <li>Your top artists and tracks</li>
                <li>Listening history analysis</li>
                <li>Personalized statistics</li>
            </ul>
        </div>
    </div>
</div>

<style>
    .login-container {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        padding: 20px;
    }

    .login-card {
        background: white;
        padding: 40px;
        border-radius: 12px;
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
        text-align: center;
        max-width: 400px;
        width: 100%;
    }

    h1 {
        color: #333;
        margin-bottom: 10px;
        font-size: 2rem;
    }

    p {
        color: #666;
        margin-bottom: 30px;
        line-height: 1.6;
    }

    .spotify-login-btn {
        background: #1db954;
        color: white;
        border: none;
        padding: 15px 30px;
        border-radius: 25px;
        font-size: 1.1rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 10px;
        width: 100%;
        transition: background-color 0.3s ease;
    }

    .spotify-login-btn:hover {
        background: #1ed760;
    }

    .features {
        margin-top: 30px;
        text-align: left;
    }

    .features h3 {
        color: #333;
        margin-bottom: 15px;
    }

    .features ul {
        list-style: none;
        padding: 0;
    }

    .features li {
        color: #666;
        padding: 8px 0;
        position: relative;
        padding-left: 20px;
    }

    .features li:before {
        content: "✓";
        color: #1db954;
        font-weight: bold;
        position: absolute;
        left: 0;
    }
</style>
