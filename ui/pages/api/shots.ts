import { withApiAuthRequired, getAccessToken } from '@auth0/nextjs-auth0';

export default withApiAuthRequired(async function shows(req, res) {
    try {
        const { accessToken } = await getAccessToken(req, res, {
        });

        let shotsApiUrl = process.env.SHOTS_API_URL;
        const response = await fetch(shotsApiUrl, {
            headers: {
                Authorization: `Bearer ${accessToken}`
            }
        });

        const shows = await response.json();
        res.status(response.status || 200).json(shows);
    } catch (error) {
        console.error(error);
        res.status(error.status || 500).json({
            code: error.code,
            error: error.message
        });
    }
});