import {withPageAuthRequired} from '@auth0/nextjs-auth0';
import useSwr from 'swr'
import React from "react";
import {Shot} from "../interfaces";
import Layout from "../components/Layout";
import ShotList from "../components/ShotList";


const Shots: React.FC = () => {

    const fetcher = (url: string) => fetch(url).then((res) => res.json())
    const {data, error} = useSwr<Array<Shot>>('/api/shots', fetcher)

    if (error) return <div>Failed to load shots</div>
    if (!data) return <div>Loading...</div>

    return (
        <Layout>
            <div>
                <h1 className="my-8 font-black tracking-tight text-xl">Latest shots:</h1>
                <ShotList items={data}></ShotList>
            </div>
        </Layout>
    )
}

export default Shots

export const getServerSideProps = withPageAuthRequired()