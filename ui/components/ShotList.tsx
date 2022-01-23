import * as React from 'react'
import {Shot} from "../interfaces";

type Props = {
  items: Shot[]
}

const ShotList = ({ items }: Props) => (
  <ul className="list-disc">
    {items.map((item) => (
      <li key={item.key}>
        {item.key}
      </li>
    ))}
  </ul>
)

export default ShotList
