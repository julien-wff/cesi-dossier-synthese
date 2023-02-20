import tabula
import pandas as pd
import re


"""
@description: Read the pdf file and return a list of dictionaries representing the grades
@param path: The path to the pdf file
@return: A list of dictionaries representing the grades
"""
def read_pdf(path):
    # Read all the tables in the pdf
    grades = tabula.read_pdf(path, pages='all')

    # Concatenate all the tables into one if they have a 'AXE / U.E. / E.E.' and 'Coeff.' column
    grades = pd.concat(
        [table for table in grades if 'AXE / U.E. / E.E.' in table.columns and 'Coeff.' in table.columns],
        ignore_index=True
    )

    # Convert the DataFrame to a dictionary
    grades = grades.to_dict(orient='records')

    result = []

    for row in grades:
        # Reformat the row
        grade = {
            'name': row['AXE / U.E. / E.E.'].replace('\r', ' ').strip(),
            'grade': None,
            'coefficient': None if pd.isna(row['Coeff.']) else int(row['Coeff.']),
        }

        if row['Evaluation'] in ('A', 'B', 'C', 'D'):
            grade['grade'] = row['Evaluation']
        elif re.match(r'[A-D]( / [A-D])?', str(row['Evaluation'])):
            grade['previous_grade'] = row['Evaluation'].split(' / ')[0]
            grade['grade'] = row['Evaluation'].split(' / ')[1]

        # Put the row in the result list, according to its type
        if grade['coefficient'] is None:
            if not grade['name'].startswith('['):
                result.append({
                    'name': grade['name'],
                    'categories': [],
                })
            else:
                result[-1]['categories'].append({
                    'name': grade['name'],
                    'grades': [],
                })
        else:
            result[-1]['categories'][-1]['grades'].append(grade)

    # Print a formatted JSON object
    return result
