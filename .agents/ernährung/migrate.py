import sqlite3
import re
from datetime import datetime

db_path = '/Users/dennisschroder/private/.agents/ernährung/nutrition.db'
daily_log_path = '/Users/dennisschroder/private/.agents/ernährung/daily_log.md'
weight_log_path = '/Users/dennisschroder/private/.agents/ernährung/weight_log.md'

def migrate_weight():
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    with open(weight_log_path, 'r') as f:
        content = f.read()
    
    matches = re.findall(r'\| (?:Mo|Di|Mi|Do|Fr|Sa|So) (\d{2}\.\d{2})\. \| (\d+\.\d+) kg \| (.*?) \|', content)
    for date_str, weight, comment in matches:
        full_date = f"2026-{date_str.split('.')[1]}-{date_str.split('.')[0]}"
        try:
            cursor.execute("INSERT OR REPLACE INTO weight_logs (date, weight, comment) VALUES (?, ?, ?)", (full_date, float(weight), comment.strip()))
        except Exception as e:
            print(f"Error inserting weight for {full_date}: {e}")
    conn.commit()
    conn.close()

def migrate_daily():
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    with open(daily_log_path, 'r') as f:
        content = f.read()
    
    days = re.split(r'## ', content)[1:]
    for day_content in days:
        lines = day_content.split('\n')
        header = lines[0]
        date_match = re.search(r'(\d{2}\.\d{2}\.\d{4})', header)
        day_match = re.search(r'\(Tag (\d+)\)', header)
        
        if not date_match: continue
        
        raw_date = date_match.group(1)
        # Convert DD.MM.YYYY to YYYY-MM-DD
        date_obj = datetime.strptime(raw_date, '%d.%m.%Y')
        formatted_date = date_obj.strftime('%Y-%m-%d')
        day_num = int(day_match.group(1)) if day_match else None
        
        # Targets
        cal_target = re.search(r'Kalorien: < (\d+)', day_content)
        prot_target = re.search(r'Protein: (\d+)', day_content)
        
        # Actuals
        totals = re.search(r'\*\*Gesamt\*\*: ~(\d+) kcal / (\d+) kcal\n\*\*Protein\*\*: ~(\d+)g / (\d+)g\n\*\*Fett\*\*: ~(\d+)g / <(\d+)g\n\*\*Carbs\*\*: ~(\d+)g / <(\d+)g', day_content)
        
        # Conclusion
        conclusion = re.search(r'\*\*Fazit\*\*: (.*)', day_content)
        
        cursor.execute("""
            INSERT OR REPLACE INTO daily_logs 
            (date, day_number, calories_target, protein_target, fat_target, carbs_target, 
             calories_actual, protein_actual, fat_actual, carbs_actual, conclusion)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """, (
            formatted_date, day_num, 
            int(cal_target.group(1)) if cal_target else 1500,
            int(prot_target.group(1)) if prot_target else 200,
            30, 50,
            int(totals.group(1)) if totals else None,
            int(totals.group(3)) if totals else None,
            int(totals.group(5)) if totals else None,
            int(totals.group(7)) if totals else None,
            conclusion.group(1) if conclusion else None
        ))
    
    conn.commit()
    conn.close()

if __name__ == "__main__":
    migrate_weight()
    migrate_daily()
    print("Migration complete.")
