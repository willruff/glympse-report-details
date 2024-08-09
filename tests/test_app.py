import unittest
from report_details.app import sortFile

class TestReportDetails(unittest.TestCase):
    def test_sortFile(self):
        cases = [{
            "input": """org_id,id,agent_id\n71236,picture,3452\n71237,flipper,3453\n71238,picture,3454""",
            "expected": (["org_id","agent_id"],
                [{
                "org_id": '71236',
                "agent_id": '3452',
                },
                {
                "org_id": '71238',
                "agent_id": '3454',
                }]), 
            "msg": "2 rows should be returned with the picture column excluded"
            },   
            {
            "input": """org_id,id,agent_id\n71236,picture,3452\n71237,picture,3453\n71238,picture,3454""", 
            "expected": (["org_id", "agent_id"],
                [{
                "org_id": '71236',
                "agent_id": '3452',
                },
                {
                "org_id": '71237',
                "agent_id": '3453',   
                },
                {
                "org_id": '71238',
                "agent_id": '3454',
                }]), 
            "msg": "All rows should be returned"
            },
            {
            "input": """org_id,id,agent_id\n71236,desk,3452\n71237,chair,3453\n71238,lamp,3454""", 
            "expected":(["org_id","agent_id"], []),
            "msg": "None of the rows match search filter"
            },
            {
            "input": """""", 
            "expected": (None, None),
            "msg": "No rows should be returned. Value should be (None, None)"
            },
            {
            "input": """org_id,flags,agent_id\n71236,tes,3452\n71237,flipper,3453\n71238,test,3454""", 
            "expected":(["org_id", "flags", "agent_id"],[]),
            "msg": "Only the fieldnames should be returned"
            },
        ]

        for case in cases:
            recipients = case["input"]
            actual = sortFile(recipients, "picture", "id")
            self.assertEqual(actual, case["expected"], case["msg"])
