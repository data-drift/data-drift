import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import * as duckdb from "@duckdb/duckdb-wasm";
import duckdb_wasm from "@duckdb/duckdb-wasm/dist/duckdb-mvp.wasm?url";
import mvp_worker from "@duckdb/duckdb-wasm/dist/duckdb-browser-mvp.worker.js?url";
import duckdb_wasm_eh from "@duckdb/duckdb-wasm/dist/duckdb-eh.wasm?url";
import eh_worker from "@duckdb/duckdb-wasm/dist/duckdb-browser-eh.worker.js?url";

interface DuckDbProviderProps {
  children: ReactNode; // Typing children using ReactNode
}

const MANUAL_BUNDLES = {
  mvp: {
    mainModule: duckdb_wasm,
    mainWorker: mvp_worker,
  },
  eh: {
    mainModule: duckdb_wasm_eh,
    mainWorker: eh_worker,
  },
} as const;

// Create a Context
const DuckDbContext = createContext<duckdb.AsyncDuckDBConnection | null>(null);

const useDuckDb = () => useContext(DuckDbContext);

const DuckDbProvider = ({ children }: DuckDbProviderProps) => {
  const [db, setDb] = useState<duckdb.AsyncDuckDBConnection | null>(null);

  useEffect(() => {
    // Initialize your DuckDb instance here.
    // This is a placeholder for actual DuckDb initialization logic.
    // You might need to asynchronously load the database or its schema.
    const initializeDuckDb = async () => {
      try {
        const bundle = await duckdb.selectBundle(MANUAL_BUNDLES);
        if (!bundle.mainWorker) {
          throw new Error("No worker found in the selected bundle");
        }
        const worker = new Worker(bundle.mainWorker);
        const logger = new duckdb.ConsoleLogger();
        const dbInstance = new duckdb.AsyncDuckDB(logger, worker);
        await dbInstance.instantiate(bundle.mainModule, bundle.pthreadWorker);
        const connection = await dbInstance.connect();

        setDb(connection);
      } catch (error) {
        console.error("Failed to initialize DuckDB:", error);
      }
    };

    void initializeDuckDb();
  }, []);

  return <DuckDbContext.Provider value={db}>{children}</DuckDbContext.Provider>;
};

DuckDbProvider.useDuckDb = useDuckDb;

export default DuckDbProvider;
