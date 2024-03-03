from neo4j import GraphDatabase


class UserInfo:

    def __init__(self, uri, user, password):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))

    def close(self):
        self.driver.close()

    def get_user_info(self, name):
        with self.driver.session() as session:
            records = session.execute_read(self._get_user_info, name)

            if records is None:
                return None

            info = records[0]

            return {
                "limit": info["limit"],
                "openAIKey": info["openAIKey"],
                "cosine": info["cosine"],
            }

    @staticmethod
    def _get_user_info(tx, name):
        result = tx.run("MATCH (u:User { name: $name })" "RETURN u", name=name)

        return result.single()
