# -*- coding: utf-8 -*-

import re
from functools import partial
from lxml import etree
from parse import fastiter
from parse import es


INDEX = 'dba'
QUESTION_TYPE = 'question'
ANSWER_TYPE = 'answer'
INPUT_POSTS = ('./stackexchange/dba.stackexchange.com/Posts.xml')
INPUT_USERS = ('./stackexchange/dba.stackexchange.com/Users.xml')
INPUT_COMMENTS = ('./stackexchange/dba.stackexchange.com/Comments.xml')


def parse_tags(input):
    p = re.compile(r'<(.*?)>')
    return p.findall(input)


# put mapping
question_mapping = {
    'properties': {
        'title': {'type': 'string'},
        'author': {
            'type': 'object',
            'properties': {
                'id': {'type': 'long'},
                'name': {'type': 'string'},
            }
        },
        'creation_date': {'type': 'date'},
        'rating': {'type': 'long'},
        'tags': {'type': 'string'},
        'body': {'type': 'string'},
        'comments': {
            'type': 'nested',
            'properties': {
                'creation_date': {'type': 'date'},
                'author': {
                    'type': 'object',
                    'properties': {
                        'id': {'type': 'long'},
                        'name': {'type': 'string'},
                    }
                },
                'rating': {'type': 'long'},
                'body': {'type': 'string'},
            }
        }
    }
}

answer_mapping = {
    '_parent': {'type': 'question'},
    'properties': {
        'author': {
            'type': 'object',
            'properties': {
                'id': {'type': 'long'},
                'name': {'type': 'string'},
            }
        },
        'creation_date': {'type': 'date'},
        'accepted': {'type': 'boolean'},
        'rating': {'type': 'long'},
        'body': {'type': 'string'},
        'comments': {
            'type': 'nested',
            'properties': {
                'creation_date': {'type': 'date'},
                'author': {
                    'type': 'object',
                    'properties': {
                        'id': {'type': 'long'},
                        'name': {'type': 'string'},
                    }
                },
                'rating': {'type': 'long'},
                'body': {'type': 'string'},
            }
        }
    }
}

es.indices.delete(index=INDEX, ignore=[404, 400])
es.indices.create(index=INDEX)
es.indices.put_mapping(index=INDEX, doc_type=QUESTION_TYPE,
                       body=question_mapping)
es.indices.put_mapping(index=INDEX, doc_type=ANSWER_TYPE,
                       body=answer_mapping)


def process_questions(answers, users, comments, elem):
    if elem.attrib['PostTypeId'] != '1':
        return
    _id = int(elem.attrib['Id'])
    body = {
        'title': elem.attrib['Title'],
        'creation_date': elem.attrib['CreationDate'],
        'rating': int(elem.attrib['Score']),
        'tags': parse_tags(elem.attrib['Tags']),
        'body': elem.attrib['Body'],
        'comments': [],
        'author': {},
    }
    if _id in comments:
        body['comments'] = comments[_id]
    if 'AcceptedAnswerId' in elem.attrib:
        answers.add(elem.attrib['AcceptedAnswerId'])
    if 'OwnerUserId' in elem.attrib:
        author_id = int(elem.attrib['OwnerUserId'])
        body['author']['id'] = author_id
        body['author']['name'] = users[author_id]
    if 'UserDisplayName' in elem.attrib:
        body['author']['name'] = elem.attrib['UserDisplayName']
    es.create(index=INDEX, doc_type=QUESTION_TYPE, id=_id, body=body)


def process_answers(answers, users, comments,  elem):
    if elem.attrib['PostTypeId'] != '2':
        return
    _id = int(elem.attrib['Id'])
    _parent = elem.attrib['ParentId']
    body = {
        'accepted': _id in answers,
        'creation_date': elem.attrib['CreationDate'],
        'rating': int(elem.attrib['Score']),
        'body': elem.attrib['Body'],
        'comments': [],
        'author': {},
    }
    if _id in comments:
        body['comments'] = comments[_id]
    if 'OwnerUserId' in elem.attrib:
        author_id = int(elem.attrib['OwnerUserId'])
        body['author']['id'] = author_id
        body['author']['name'] = users[author_id]
    if 'UserDisplayName' in elem.attrib:
        body['author']['name'] = elem.attrib['UserDisplayName']
    es.create(index=INDEX, doc_type=ANSWER_TYPE, id=_id, body=body,
              parent=_parent)


def process_users(collector, elem):
    collector[int(elem.attrib['Id'])] = elem.attrib['DisplayName']


def process_comments(collector, users, elem):
    post_id = int(elem.attrib['PostId'])
    if post_id not in collector:
        collector[post_id] = []
    body = {
        'creation_date': elem.attrib['CreationDate'],
        'rating': int(elem.attrib['Score']),
        'body': elem.attrib['Text'],
        'author': {},
    }
    if 'UserId' in elem.attrib:
        author_id = int(elem.attrib['UserId'])
        body['author']['id'] = author_id
        body['author']['name'] = users[author_id]
    if 'UserDisplayName' in elem.attrib:
        body['author']['name'] = elem.attrib['UserDisplayName']
    collector[post_id].append(body)


context = etree.iterparse(INPUT_USERS, tag='row')
users = {}
fastiter(context, partial(process_users, users))

context = etree.iterparse(INPUT_COMMENTS, tag='row')
comments = {}
fastiter(context, partial(process_comments, comments, users))

answers = set()
context = etree.iterparse(INPUT_POSTS, tag='row')
fastiter(context, partial(process_questions, answers, users, comments))

context = etree.iterparse(INPUT_POSTS, tag='row')
fastiter(context, partial(process_answers, answers, users, comments))
