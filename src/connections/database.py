import os

from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import as_declarative
from sqlalchemy.orm import scoped_session, sessionmaker

_session = None

Session = scoped_session(
    sessionmaker(autocommit=False, autoflush=False),
)


@as_declarative()
class Base(object):
    pass


def bind_engine():
    engine = create_engine(os.getenv('DATABASE_URL'))
    Base.metadata.bind = engine


def get_session():
    global _session
    if _session:
        return _session

    _session = bind_engine()

    return _session


def get_db():
    db = Session()
    try:
        yield db
    finally:
        db.close()


def create_db():
    Base.metadata.create_all(bind_engine())
