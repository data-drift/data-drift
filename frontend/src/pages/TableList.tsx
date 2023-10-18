import { useLoaderData } from "react-router-dom";
import { getTableList } from "../services/data-drift";

const loader = async () => {
  const tableList = await getTableList();
  return tableList;
};

type LoaderData = Awaited<ReturnType<typeof loader>>;

const TableList = () => {
  const loader = useLoaderData() as LoaderData;
  return (
    <div>
      <h1>Store: {loader.data.store}</h1>
      <ul>
        {loader.data.tables.map((table) => (
          <li key={table} style={{ textAlign: "justify" }}>
            <a href={`./tables/${table}`}>{table}</a>
          </li>
        ))}
      </ul>
    </div>
  );
};

TableList.loader = loader;

export default TableList;
