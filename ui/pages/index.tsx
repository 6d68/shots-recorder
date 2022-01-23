import React from "react";
import Layout from "../components/Layout";
import {useUser} from "@auth0/nextjs-auth0";
import {useRouter} from 'next/router'


const Homepage: React.FC = () => {

    const {user} = useUser();
    const router = useRouter()

    if (user) {
        router.replace({
            pathname: '/shots'
        });
    }

    return (
        <Layout>
            <div>
                <h1 className="text-center my-8 font-black tracking-tight text-6xl">Shots Recorder</h1>
                <h2 className="text-center">Login to view latest shots</h2>
            </div>
        </Layout>
    )
}

export default Homepage;
